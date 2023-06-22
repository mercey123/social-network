package sqlite

import (
	"database/sql"
	eh "social-network/packages/errorHandler"
	"social-network/packages/models"
	"time"
)

func CreateGroup(group models.CreateGroupRequest, userId int) (*models.CreateGroupResponse, error) {
	sqlStmt, err := db.Prepare(`INSERT INTO groups (userId, title, description, creationDate)
	VALUES (?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}

	creationDate := time.Now().UnixMilli()

	result, err := sqlStmt.Exec(userId, group.Title, group.Description, creationDate)
	if err != nil {
		return nil, err
	}

	groupId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.CreateGroupResponse{
		Id:           int(groupId),
		Title:        group.Title,
		Description:  group.Description,
		CreationDate: int(creationDate),
		UserId:       userId,
	}, nil
}

func GetAllGroups(userId int) ([]models.GetGroupInfoResponse, error) {
	groups := make([]models.GetGroupInfoResponse, 0)

	rows, err := db.Query(`
	SELECT id,
		title,
		userId,
		(
			SELECT isAccepted
			FROM group_members
			WHERE id = group_members.groupId
				AND group_members.userId = ?
		),
		(
			SELECT 1
			FROM group_invitation
			WHERE id = group_invitation.groupId
				AND group_invitation.userId = ?
		)
	FROM groups`, userId, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tempGroup := models.GetGroupInfoResponse{JoinStatus: 1}
		ownerId := 0
		var isAccepted *int
		var isInvited *int

		err = rows.Scan(&tempGroup.Id, &tempGroup.Title, &ownerId, &isAccepted, &isInvited)
		if err != nil {
			return nil, err
		}

		switch {
		case ownerId == userId || (isAccepted != nil && *isAccepted == 1):
			tempGroup.JoinStatus = 3
		case isAccepted != nil && *isAccepted == 0:
			tempGroup.JoinStatus = 2
		case isInvited != nil:
			tempGroup.JoinStatus = 4
		default:
			tempGroup.JoinStatus = 1
		}

		groups = append(groups, tempGroup)
	}

	return groups, nil
}

func GetGroupById(groupId, userId int) (*models.GetGroupResponse, error) {
	group := &models.GetGroupResponse{JoinStatus: 1}
	ownerId := 0
	var isAccepted *int
	var isInvited *int

	err := db.QueryRow(`
	SELECT id,
		title,
		description,
		userId,
		(
			SELECT isAccepted
			FROM group_members
			WHERE id = group_members.groupId
				AND group_members.userId = ?
		),
		(
			SELECT 1
			FROM group_invitation
			WHERE id = group_invitation.groupId
				AND group_invitation.userId = ?
		)
	FROM groups
	WHERE id = ?`, userId, userId, groupId).Scan(&group.Id, &group.Title, &group.Description, &ownerId, &isAccepted, &isInvited)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, eh.NewErrorResponse(eh.ErrNotFound, "wrong variable(s) in request")
	}

	switch {
	case ownerId == userId || (isAccepted != nil && *isAccepted == 1):
		group.JoinStatus = 3
	case isAccepted != nil && *isAccepted == 0:
		group.JoinStatus = 2
	case isInvited != nil:
		group.JoinStatus = 4
	default:
		group.JoinStatus = 1
	}

	return group, nil
}

func InviteToGroup(groupId, userId int, invitedUsers []int) ([]int, error) {
	sqlStmt, err := db.Prepare(`INSERT INTO group_invitation (groupId, userId)
	VALUES (?, ?);`)
	if err != nil {
		return nil, err
	}

	group, err := GetGroupById(groupId, userId)
	if err != nil {
		return nil, err
	}

	if group.JoinStatus != 3 {
		return nil, eh.NewErrorResponse(eh.ErrNoAccess, "no access to this action")
	}

	temp := make([]int, 0)

	for _, id := range invitedUsers {
		group, err = GetGroupById(groupId, id)
		if group.JoinStatus != 1 {
			continue
		}

		if err != nil {
			return nil, err
		}

		_, err = sqlStmt.Exec(groupId, id)
		if err != nil {
			return nil, err
		}

		temp = append(temp, id)
	}

	return temp, nil
}

func UpdateJoinStatus(groupId, userId int, isJoining bool) (int, error) {
	group, err := GetGroupById(groupId, userId)
	if err != nil {
		return 0, err
	}

	if group.JoinStatus != 1 && isJoining {
		return group.JoinStatus, nil
	} else if group.JoinStatus == 1 && !isJoining {
		return 1, nil
	}

	var sqlStmt *sql.Stmt
	if isJoining {
		sqlStmt, err = db.Prepare("INSERT INTO group_members (groupId, userId) VALUES (?, ?)")
	} else {
		sqlStmt, err = db.Prepare("DELETE FROM group_members WHERE groupId = ? AND userId = ?")
	}
	if err != nil {
		return 0, err
	}

	_, err = sqlStmt.Exec(groupId, userId)
	if err != nil {
		return 0, err
	}

	if isJoining {
		return 2, nil
	} else {
		return 1, nil
	}
}
