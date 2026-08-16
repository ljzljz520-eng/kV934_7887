const app = document.querySelector("#app");
const params = new URLSearchParams(window.location.search);
const slug = params.get("course") || "go-from-zero";
const destination = params.get("destination");
const adminMode = params.get("admin") === "1";
let visitId = "";

const icons = { questionnaire: "01", drive: "02", community: "03" };
const destinationCopy = {
  questionnaire: ["问卷已打开", "感谢你留下学习目标，提交后返回试听页继续。"],
  drive: ["资料入口", "课程资料页已准备好，可返回试听页继续访问下一步。"],
  community: ["社群入口", "欢迎加入学习社群，返回试听页完成访问链路。"]
};

function escapeHTML(value) {
  return String(value).replace(/[&<>\"]/g, (character) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", "\"": "&quot;" }[character]));
}

function layout(content, navigation = `<a class="topbar-status" href="/?admin=1&course=${encodeURIComponent(slug)}">管理</a>`) {
  app.innerHTML = `<div class="topbar"><span class="brand-mark">课</span><span>课程试听资料页</span>${navigation}</div><section class="content">${content}</section>`;
}

function renderDestination(kind) {
  const copy = destinationCopy[kind] || ["目标页面", "这是试听链路中的一个本地目标页。"];
  layout(`<div class="destination-page"><span class="eyebrow">COURSE TRIAL / DESTINATION</span><div class="destination-icon">${icons[kind] || "→"}</div><h1>${copy[0]}</h1><p>${copy[1]}</p><a class="button secondary" href="/?course=${encodeURIComponent(slug)}">返回试听页</a></div>`);
}

function renderPage(page) {
  layout(`<div class="hero"><div class="hero-image"><img src="/assets/course-workspace.jpg" alt="课程录制工作台" /></div><div class="hero-shade"></div><div class="hero-copy"><span class="eyebrow">COURSE TRIAL / ${escapeHTML(page.slug)}</span><h1>${escapeHTML(page.title)}</h1><p class="intro">${escapeHTML(page.introduction)}</p><div class="video-note"><span class="video-dot"></span><span>${escapeHTML(page.video_note)}</span></div></div></div><section class="trial-panel"><div class="panel-heading"><div><span class="eyebrow">VISITOR PATH</span><h2>按顺序完成试听访问</h2></div><span class="count-label">已开放 ${page.access_count} / ${page.access_limit || "∞"} 次</span></div><div id="path" class="path"><span>01 问卷</span><i></i><span>02 资料</span><i></i><span>03 社群</span></div><div id="action"></div></section>`);
  const action = document.querySelector("#action");
  if (!page.available) {
    action.innerHTML = `<div class="closed"><span class="closed-mark">×</span><div><strong>本次试听已结束</strong><p>${escapeHTML(page.closing_copy)}</p></div></div>`;
    return;
  }
  action.innerHTML = `<button class="button primary" id="start">${escapeHTML(page.button_label)}<span>→</span></button>`;
  document.querySelector("#start").addEventListener("click", startVisit);
}

function renderAdmin(page, notice = "") {
  const entryRows = page.entries.map((entry, index) => `<fieldset class="entry-editor" data-entry="${index}"><legend>${icons[entry.kind] || "--"} ${escapeHTML(entry.kind)}</legend><div class="field-grid"><label><span>按钮名称</span><input name="entry-label-${index}" value="${escapeHTML(entry.label)}" required /></label><label><span>目标地址</span><input name="entry-url-${index}" value="${escapeHTML(entry.url)}" required /></label></div><label><span>入口说明</span><input name="entry-summary-${index}" value="${escapeHTML(entry.summary)}" /></label><label class="toggle"><input type="checkbox" name="entry-enabled-${index}" ${entry.enabled ? "checked" : ""} /><span>启用入口</span></label></fieldset>`).join("");
  layout(`<section class="admin"><div class="admin-heading"><div><span class="eyebrow">COURSE ADMIN</span><h1>编辑试听资料</h1><p>访问次数 ${page.access_count}，页面状态 ${page.active ? "开放" : "停用"}</p></div><span id="save-status" class="save-status">${escapeHTML(notice)}</span></div><form id="admin-form"><section class="form-section"><h2>页面内容</h2><label><span>课程名称</span><input name="title" value="${escapeHTML(page.title)}" required /></label><label><span>课程介绍</span><textarea name="introduction" rows="4" required>${escapeHTML(page.introduction)}</textarea></label><label><span>试听视频说明</span><textarea name="video-note" rows="3" required>${escapeHTML(page.video_note)}</textarea></label><div class="field-grid"><label><span>主按钮文案</span><input name="button-label" value="${escapeHTML(page.button_label)}" required /></label><label><span>访问次数上限</span><input name="access-limit" type="number" min="0" value="${page.access_limit}" required /></label></div><label><span>自定义收尾文案</span><textarea name="closing-copy" rows="3" required>${escapeHTML(page.closing_copy)}</textarea></label></section><section class="form-section"><h2>访问入口</h2>${entryRows}</section><div class="form-actions"><button class="button primary" type="submit">保存课程</button></div></form></section>`, `<a class="topbar-status" href="/?course=${encodeURIComponent(slug)}">返回访客页</a>`);
  document.querySelector("#admin-form").addEventListener("submit", (event) => saveAdmin(event, page));
}

async function saveAdmin(event, page) {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  const update = {
    title: form.get("title"),
    introduction: form.get("introduction"),
    video_note: form.get("video-note"),
    button_label: form.get("button-label"),
    closing_copy: form.get("closing-copy"),
    access_limit: Number(form.get("access-limit")),
    entries: page.entries.map((entry, index) => ({ ...entry, label: form.get(`entry-label-${index}`), summary: form.get(`entry-summary-${index}`), url: form.get(`entry-url-${index}`), enabled: form.has(`entry-enabled-${index}`) }))
  };
  const response = await fetch(`/api/admin/pages/${encodeURIComponent(slug)}`, { method: "PUT", headers: { "Content-Type": "application/json" }, body: JSON.stringify(update) });
  if (!response.ok) {
    const error = await response.json();
    document.querySelector("#save-status").textContent = error.message;
    return;
  }
  renderAdmin(await response.json(), "已保存");
}

async function startVisit() {
  const response = await fetch(`/api/pages/${encodeURIComponent(slug)}/visits`, { method: "POST" });
  if (!response.ok) {
    const error = await response.json();
    renderPage({ title: "", slug, introduction: "", video_note: "", access_count: 0, access_limit: 0, available: false, closing_copy: error.closing_copy || error.message });
    return;
  }
  const started = await response.json();
  visitId = started.visit_id;
  await nextStep();
}

async function nextStep() {
  const response = await fetch(`/api/visits/${encodeURIComponent(visitId)}/next`, { method: "POST" });
  const step = await response.json();
  const action = document.querySelector("#action");
  if (step.done) {
    action.innerHTML = `<div class="closed complete"><span class="closed-mark">✓</span><div><strong>访问链路已完成</strong><p>${escapeHTML(step.closing_copy)}</p></div></div>`;
    return;
  }
  action.innerHTML = `<div class="step-card"><div class="step-number">${String(step.step).padStart(2, "0")}</div><div class="step-copy"><span class="eyebrow">NEXT DESTINATION</span><h3>${escapeHTML(step.label)}</h3><p>${escapeHTML(step.summary)}</p></div><a class="button secondary" href="${escapeHTML(step.url)}" target="_blank" rel="noreferrer">打开入口 <span>↗</span></a><button class="button text-button" id="continue">继续 <span>→</span></button></div>`;
  document.querySelector("#continue").addEventListener("click", nextStep);
}

async function load() {
  if (destination) {
    renderDestination(destination);
    return;
  }
  try {
    if (adminMode) {
      const response = await fetch(`/api/admin/pages/${encodeURIComponent(slug)}`);
      if (!response.ok) throw new Error("课程配置不存在");
      renderAdmin(await response.json());
      return;
    }
    const response = await fetch(`/api/pages/${encodeURIComponent(slug)}`);
    if (!response.ok) throw new Error("课程资料页不存在");
    renderPage(await response.json());
  } catch (error) {
    layout(`<div class="destination-page"><span class="eyebrow">COURSE TRIAL</span><h1>暂时无法打开</h1><p>${escapeHTML(error.message)}</p></div>`);
  }
}

load();
