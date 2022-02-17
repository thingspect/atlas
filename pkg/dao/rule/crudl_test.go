//go:build !unit

package rule

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("dao-rule", createOrg.Id)
		createRule, _ := proto.Clone(rule).(*api.Rule)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, createRule)
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.NoError(t, err)
		require.NotEqual(t, rule.Id, createRule.Id)
		require.WithinDuration(t, time.Now(), createRule.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createRule.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("dao-rule", createOrg.Id)
		rule.Attr = "dao-rule-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, rule)
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.Nil(t, createRule)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Read rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readRule, err := globalRuleDAO.Read(ctx, createRule.Id,
			createRule.OrgId)
		t.Logf("readRule, err: %+v, %v", readRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule, readRule)
	})

	t.Run("Read rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readRule, err := globalRuleDAO.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readRule, err: %+v, %v", readRule, err)
		require.Nil(t, readRule)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readRule, err := globalRuleDAO.Read(ctx, createRule.Id,
			uuid.NewString())
		t.Logf("readRule, err: %+v, %v", readRule, err)
		require.Nil(t, readRule)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read rule by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readRule, err := globalRuleDAO.Read(ctx, random.String(10),
			createRule.OrgId)
		t.Logf("readRule, err: %+v, %v", readRule, err)
		require.Nil(t, readRule)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Name = "dao-rule-" + random.String(10)
		createRule.Status = api.Status_DISABLED
		updateRule, _ := proto.Clone(createRule).(*api.Rule)

		updateRule, err = globalRuleDAO.Update(ctx, updateRule)
		t.Logf("createRule, updateRule, err: %+v, %+v, %v", createRule,
			updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.Name, updateRule.Name)
		require.Equal(t, createRule.Status, updateRule.Status)
		require.True(t, updateRule.UpdatedAt.AsTime().After(
			updateRule.CreatedAt.AsTime()))
		require.WithinDuration(t, createRule.CreatedAt.AsTime(),
			updateRule.UpdatedAt.AsTime(), 2*time.Second)

		readRule, err := globalRuleDAO.Read(ctx, createRule.Id,
			createRule.OrgId)
		t.Logf("readRule, err: %+v, %v", readRule, err)
		require.NoError(t, err)
		require.Equal(t, updateRule, readRule)
	})

	t.Run("Update unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateRule, err := globalRuleDAO.Update(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.OrgId = uuid.NewString()
		createRule.Name = "dao-rule-" + random.String(10)

		updateRule, err := globalRuleDAO.Update(ctx, createRule)
		t.Logf("createRule, updateRule, err: %+v, %+v, %v", createRule,
			updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update rule by invalid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Attr = "dao-rule-" + random.String(40)

		updateRule, err := globalRuleDAO.Update(ctx, createRule)
		t.Logf("createRule, updateRule, err: %+v, %+v, %v", createRule,
			updateRule, err)
		require.Nil(t, updateRule)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		err = globalRuleDAO.Delete(ctx, createRule.Id, createOrg.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read rule by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readRule, err := globalRuleDAO.Read(ctx, createRule.Id,
				createOrg.Id)
			t.Logf("readRule, err: %+v, %v", readRule, err)
			require.Nil(t, readRule)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalRuleDAO.Delete(ctx, uuid.NewString(), createOrg.Id)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		err = globalRuleDAO.Delete(ctx, createRule.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	ruleIDs := []string{}
	ruleNames := []string{}
	ruleStatuses := []api.Status{}
	ruleTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-rule",
			createOrg.Id))
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		ruleIDs = append(ruleIDs, createRule.Id)
		ruleNames = append(ruleNames, createRule.Name)
		ruleStatuses = append(ruleStatuses, createRule.Status)
		ruleTSes = append(ruleTSes, createRule.CreatedAt.AsTime())
	}

	t.Run("List rules by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, listCount, err := globalRuleDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0)
		t.Logf("listRules, listCount, err: %+v, %v, %v", listRules, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listRules, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, rule := range listRules {
			if rule.Id == ruleIDs[len(ruleIDs)-1] &&
				rule.Name == ruleNames[len(ruleNames)-1] &&
				rule.Status == ruleStatuses[len(ruleStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List rules by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, listCount, err := globalRuleDAO.List(ctx, createOrg.Id,
			ruleTSes[0], ruleIDs[0], 5)
		t.Logf("listRules, listCount, err: %+v, %v, %v", listRules, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listRules, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, rule := range listRules {
			if rule.Id == ruleIDs[len(ruleIDs)-1] &&
				rule.Name == ruleNames[len(ruleNames)-1] &&
				rule.Status == ruleStatuses[len(ruleStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List rules by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, listCount, err := globalRuleDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 1)
		t.Logf("listRules, listCount, err: %+v, %v, %v", listRules, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listRules, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, listCount, err := globalRuleDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0)
		t.Logf("listRules, listCount, err: %+v, %v, %v", listRules, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listRules, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List rules by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, listCount, err := globalRuleDAO.List(ctx, random.String(10),
			time.Time{}, "", 0)
		t.Logf("listRules, listCount, err: %+v, %v, %v", listRules, listCount,
			err)
		require.Nil(t, listRules)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestListByTags(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-rule"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	ruleIDs := []string{}
	ruleDeviceTags := []string{}
	ruleAttrs := []string{}
	for i := 0; i < 3; i++ {
		rule := random.Rule("dao-rule", createOrg.Id)
		rule.Status = api.Status_ACTIVE
		createRule, err := globalRuleDAO.Create(ctx, rule)
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		ruleIDs = append(ruleIDs, createRule.Id)
		ruleDeviceTags = append(ruleDeviceTags, createRule.DeviceTag)
		ruleAttrs = append(ruleAttrs, createRule.Attr)
	}

	t.Run("List rules by valid org ID and unique attr", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, err := globalRuleDAO.ListByTags(ctx, createOrg.Id,
			ruleAttrs[len(ruleAttrs)-1], ruleDeviceTags)
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules, 1)
		require.Equal(t, listRules[0].Id, ruleIDs[len(ruleIDs)-1])
		require.Equal(t, listRules[0].DeviceTag,
			ruleDeviceTags[len(ruleDeviceTags)-1])
		require.Equal(t, listRules[0].Attr, ruleAttrs[len(ruleAttrs)-1])
	})

	t.Run("List rules by valid org ID and api attr", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("dao-rule", createOrg.Id)
		rule.Status = api.Status_ACTIVE
		rule.Attr = ruleAttrs[0]

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createRule, err := globalRuleDAO.Create(ctx, rule)
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		lRuleDeviceTags := ruleDeviceTags
		lRuleDeviceTags = append(lRuleDeviceTags, createRule.DeviceTag)

		listRules, err := globalRuleDAO.ListByTags(ctx, createOrg.Id,
			ruleAttrs[0], lRuleDeviceTags)
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules, 2)
		for _, rule := range listRules {
			require.Contains(t, lRuleDeviceTags, rule.DeviceTag)
			require.Equal(t, ruleAttrs[0], rule.Attr)
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, err := globalRuleDAO.ListByTags(ctx, uuid.NewString(),
			ruleAttrs[len(ruleAttrs)-1], ruleDeviceTags)
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules, 0)
	})

	t.Run("List rules by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listRules, err := globalRuleDAO.ListByTags(ctx, random.String(10),
			ruleAttrs[len(ruleAttrs)-1], ruleDeviceTags)
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
