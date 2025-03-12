// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.8.1

package socialmodels

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	groupMembersFieldNames          = builder.RawFieldNames(&GroupMembers{})
	groupMembersRows                = strings.Join(groupMembersFieldNames, ",")
	groupMembersRowsExpectAutoSet   = strings.Join(stringx.Remove(groupMembersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	groupMembersRowsWithPlaceHolder = strings.Join(stringx.Remove(groupMembersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheGroupMembersIdPrefix = "cache:groupMembers:id:"
)

type (
	groupMembersModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*GroupMembers, error)
		Update(ctx context.Context, data *GroupMembers) error
		Delete(ctx context.Context, id uint64) error
		ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error)
		ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error)
		FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMembers, error)
	}

	defaultGroupMembersModel struct {
		sqlc.CachedConn
		table string
	}

	GroupMembers struct {
		Id          uint64         `db:"id"`
		GroupId     string         `db:"group_id"`
		UserId      string         `db:"user_id"`
		RoleLevel   int64          `db:"role_level"`
		JoinTime    sql.NullTime   `db:"join_time"`
		JoinSource  sql.NullInt64  `db:"join_source"`
		InviterUid  sql.NullString `db:"inviter_uid"`
		OperatorUid sql.NullString `db:"operator_uid"`
	}
)

func newGroupMembersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultGroupMembersModel {
	return &defaultGroupMembersModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`group_members`",
	}
}

func (m *defaultGroupMembersModel) Delete(ctx context.Context, id uint64) error {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, groupMembersIdKey)
	return err
}

func (m *defaultGroupMembersModel) FindOne(ctx context.Context, id uint64) (*GroupMembers, error) {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, id)
	var resp GroupMembers
	err := m.QueryRowCtx(ctx, &resp, groupMembersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMembersRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultGroupMembersModel) Insert(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error) {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, groupMembersRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid)
		}
		return conn.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid)
	}, groupMembersIdKey)
	return ret, err
}
func (m *defaultGroupMembersModel) ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupMembersRows, m.table)
	var resp []*GroupMembers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
func (m *defaultGroupMembersModel) ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupMembersRows, m.table)
	var resp []*GroupMembers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
func (m *defaultGroupMembersModel) Update(ctx context.Context, data *GroupMembers) error {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupMembersRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource, data.InviterUid, data.OperatorUid, data.Id)
	}, groupMembersIdKey)
	return err
}

func (m *defaultGroupMembersModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, primary)
}

func (m *defaultGroupMembersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMembersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultGroupMembersModel) tableName() string {
	return m.table
}
func (m *defaultGroupMembersModel) FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `group_id` = ?", groupMembersRows, m.table)
	var resp GroupMembers
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, groupId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}

}
