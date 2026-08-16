package trial

func FixturePages() []Page {
	return []Page{
		{
			ID: "page-go-basics", Slug: "go-from-zero", Title: "Go 从零到上线",
			Introduction: "用一条完整的项目路径，带你从第一个 Go 文件走到可发布服务。",
			VideoNote:    "试听视频覆盖模块、并发和部署的核心方法，时长约 18 分钟。",
			ButtonLabel:  "开始试听", ClosingCopy: "本期试听资料已领取完毕，感谢关注。", AccessLimit: 20, Active: true,
			Entries: []Entry{
				{Kind: EntryQuestionnaire, Label: "填写课前问卷", Summary: "留下你的基础和学习目标，老师会据此准备课程。", URL: "/?destination=questionnaire", Enabled: true},
				{Kind: EntryDrive, Label: "领取课程资料", Summary: "获取试听讲义、示例代码与课后练习。", URL: "/?destination=drive", Enabled: true},
				{Kind: EntryCommunity, Label: "加入学习社群", Summary: "进入社群和同学交流学习进度。", URL: "/?destination=community", Enabled: true},
			},
		},
		{
			ID: "page-disabled-drive", Slug: "go-disabled-drive", Title: "Go 服务实战试听",
			Introduction: "一页串起服务开发前的准备工作和学习入口。",
			VideoNote:    "试听视频说明：从接口设计开始，快速建立服务骨架。",
			ButtonLabel:  "进入试听", ClosingCopy: "本期试听资料已领取完毕，感谢关注。", AccessLimit: 20, Active: true,
			Entries: []Entry{
				{Kind: EntryQuestionnaire, Label: "填写课前问卷", Summary: "告诉老师你的服务开发经验。", URL: "/?destination=questionnaire", Enabled: true},
				{Kind: EntryCommunity, Label: "加入学习社群", Summary: "和同学一起完成试听任务。", URL: "/?destination=community", Enabled: true},
				{Kind: EntryDrive, Label: "领取课程资料", Summary: "该资料入口已暂停开放。", URL: "/?destination=drive", Enabled: false},
			},
		},
		{
			ID: "page-limited", Slug: "go-limited", Title: "Go 小班试听",
			Introduction: "固定名额的小班试听资料页。", VideoNote: "试听视频说明：小班课程重点演示。",
			ButtonLabel: "预约试听", ClosingCopy: "本期试听名额已满，感谢你的关注。", AccessLimit: 1, Active: true,
			Entries: []Entry{{Kind: EntryQuestionnaire, Label: "填写预约问卷", Summary: "提交预约信息。", URL: "/?destination=questionnaire", Enabled: true}},
		},
		{
			ID: "page-inactive", Slug: "go-archived", Title: "已归档课程", Introduction: "课程资料已归档。",
			VideoNote: "试听视频已下线。", ButtonLabel: "查看课程", ClosingCopy: "该课程已结束。", AccessLimit: 0, Active: false,
			Entries: []Entry{{Kind: EntryQuestionnaire, Label: "课前问卷", Summary: "", URL: "/?destination=questionnaire", Enabled: true}},
		},
	}
}
