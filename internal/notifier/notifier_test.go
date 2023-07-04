package notifier

import (
	"context"
	"testing"
)

func TestNotifier_Send(t *testing.T) {
	type fields struct {
		clients map[NotificationType]Client
	}

	type args struct {
		ctx              context.Context
		notification     Notification
		notificationType NotificationType
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should succeed with notifier client",
			fields: fields{
				clients: map[NotificationType]Client{
					EmailNotificationType: &mockNotifier{},
				},
			},
			args: args{
				ctx: context.Background(),
				notification: Notification{
					Message:   "hello",
					Recipient: "recipient",
				},
				notificationType: EmailNotificationType,
			},
			wantErr: false,
		},
		{
			name: "should fail with unknown notifier client",
			fields: fields{
				clients: map[NotificationType]Client{},
			},
			args: args{
				ctx: context.Background(),
				notification: Notification{
					Message:   "hello",
					Recipient: "recipient",
				},
				notificationType: EmailNotificationType,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Notifier{
				clients: tt.fields.clients,
			}
			if err := n.Send(tt.args.ctx, tt.args.notification, tt.args.notificationType); (err != nil) != tt.wantErr {
				t.Errorf("Notifier.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockNotifier struct {
}

func (m *mockNotifier) Send(ctx context.Context, notification Notification) error {
	return nil
}
