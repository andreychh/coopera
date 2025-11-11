package dialogues

import (
	"context"
)

type Topic string

const TopicNeutral Topic = "neutral"

// Dialogues is External, shared resource with a lifetime longer than a single
// method call. It holds state for all users and all requests.
type Dialogues interface {
	StartDialogue(ctx context.Context, id int64, topic Topic) (Dialogue, error)
	Dialogue(id int64) Dialogue
}

type Dialogue interface {
	ChangeTopic(ctx context.Context, topic Topic) error
	Topic(ctx context.Context) (Topic, error)
	Exists(ctx context.Context) (bool, error)
}
