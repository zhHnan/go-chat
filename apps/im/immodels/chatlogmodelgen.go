// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.1

package immodels

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DefaultChatLogCount int64 = 100

type chatLogModel interface {
	Insert(ctx context.Context, data *ChatLog) error
	FindOne(ctx context.Context, id string) (*ChatLog, error)
	Update(ctx context.Context, data *ChatLog) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id string) (int64, error)
	ListBySendTime(ctx context.Context, conversationId string, startTime, endTime, count int64) ([]*ChatLog, error)
}

type defaultChatLogModel struct {
	conn *mon.Model
}

func newDefaultChatLogModel(conn *mon.Model) *defaultChatLogModel {
	return &defaultChatLogModel{conn: conn}
}

func (m *defaultChatLogModel) Insert(ctx context.Context, data *ChatLog) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	_, err := m.conn.InsertOne(ctx, data)
	return err
}

func (m *defaultChatLogModel) FindOne(ctx context.Context, id string) (*ChatLog, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data ChatLog

	err = m.conn.FindOne(ctx, &data, bson.M{"_id": oid})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultChatLogModel) Update(ctx context.Context, data *ChatLog) (*mongo.UpdateResult, error) {
	data.UpdateAt = time.Now()

	res, err := m.conn.UpdateOne(ctx, bson.M{"_id": data.ID}, bson.M{"$set": data})
	return res, err
}

func (m *defaultChatLogModel) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrInvalidObjectId
	}

	res, err := m.conn.DeleteOne(ctx, bson.M{"_id": oid})
	return res, err
}

// ListBySendTime 根据发送时间获取聊天记录列表。
// 该方法根据指定的会话ID和时间范围，从数据库中检索聊天记录。
// 参数:
//
//	ctx - 上下文，用于取消请求和传递请求级值。
//	conversationId - 会话ID，用于标识特定的聊天会话。
//	startTime - 开始时间，表示检索的聊天记录的最早发送时间。
//	endTime - 结束时间，表示检索的聊天记录的最晚发送时间。
//	count - 希望检索的聊天记录的数量，如果为0，则使用默认值。
//
// 返回值:
//
//	[]*ChatLog - 聊天记录列表，如果找不到则为nil。
//	error - 错误信息，如果发生错误则返回。
func (m *defaultChatLogModel) ListBySendTime(ctx context.Context, conversationId string, startTime, endTime, count int64) ([]*ChatLog, error) {
	var data []*ChatLog

	// 设置查询选项，包括限制结果数量和按发送时间降序排序。
	opt := options.FindOptions{
		Limit: &DefaultChatLogCount,
		Sort: bson.M{
			"sendTime": -1,
		},
	}

	// 构建查询过滤条件，根据会话ID和发送时间范围进行筛选。
	filter := bson.M{
		"conversationId": conversationId,
	}
	if endTime > 0 {
		filter["sendTime"] = bson.M{
			"$gt":  endTime,
			"$lte": startTime,
		}
	} else {
		filter["sendTime"] = bson.M{
			"$lt": startTime,
		}
	}

	// 如果请求的数量大于0，则设置查询的限制数量。
	if count > 0 {
		opt.Limit = &count
	}

	// 执行数据库查询操作，并根据返回的错误类型进行处理。
	err := m.conn.Find(ctx, &data, filter, &opt)
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
