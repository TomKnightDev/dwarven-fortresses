package components

type Text struct {
	Content string
}

func NewText(content string) Text {
	return Text{content}
}
