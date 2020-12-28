// Package team1 contains code for team 1's client implementation
package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

// ForageOutcome is how much we put in to a foraging trip and how much we got
// back
type ForageOutcome struct {
	contribution shared.Resources
	revenue      shared.Resources
}
type ForageHistory map[shared.ForageType][]ForageOutcome

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

type clientConfig struct {
	randomForageTurns    uint
	desperationThreshold shared.Resources
	evadeTaxes           bool
}

type client struct {
	*baseclient.BaseClient

	forageHistory ForageHistory
	taxAmount     shared.Resources

	config clientConfig
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		config: clientConfig{
			randomForageTurns:    5,
			desperationThreshold: 20,
			evadeTaxes:           false,
		},
	}
}

func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

/********************/
/***  Messages      */
/********************/

func (c *client) ReceiveCommunication(
	sender shared.ClientID,
	data map[shared.CommunicationFieldName]shared.CommunicationContent,
) {
	for field, content := range data {
		switch field {
		case shared.TaxAmount:
			// if !__isPresident__(sender) {
			// }
			c.taxAmount = shared.Resources(content.IntegerData)
		}
	}
}

func (c *client) GetTaxContribution() shared.Resources {
	if !c.config.evadeTaxes {
		return 0
	}
	return c.taxAmount
}

/********************/
/***    Foraging    */
/********************/

func (c *client) RandomForage() shared.ForageDecision {
	// Up to 10% of our current resources
	forageContribution := shared.Resources(0.1*rand.Float64()) * c.gameState().ClientInfo.Resources
	c.Logf("[Forage][Decision]:Random")
	// TODO Add fish foraging when it's done
	return shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: forageContribution,
	}
}

func (c *client) NormalForage() shared.ForageDecision {
	// Find the forageType with the best average returns
	bestForageType := shared.ForageType(-1)
	bestROI := 0.0

	// TODO: We could also do this in one go. So just repeat the exact decision
	// that gave the best ROI, instead of picking forageType based on its
	// average ROI

	for forageType, outcomes := range c.forageHistory {
		ROIsum := 0.0
		for _, outcome := range outcomes {
			if outcome.contribution != 0 {
				ROIsum += float64(outcome.revenue/outcome.contribution) - 1
			}
		}

		averageROI := ROIsum / float64(len(outcomes))

		if averageROI > bestROI {
			bestROI = averageROI
			bestForageType = forageType
		}
	}

	// Not foraging is best
	if bestForageType == shared.ForageType(-1) {
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}
	}

	// Find the value of resources that gave us the best return and add some
	// noise to it. Cap to 20% of our stockpile
	pastOutcomes := c.forageHistory[bestForageType]
	bestValue := shared.Resources(0)
	bestValueROI := shared.Resources(0)
	for _, outcome := range pastOutcomes {
		if outcome.contribution != 0 {
			ROI := (outcome.revenue / outcome.contribution) - 1
			if ROI > bestValueROI {
				bestValue = outcome.contribution
				bestValueROI = ROI
			}
		}
	}
	bestValue = shared.Resources(math.Min(
		float64(bestValue),
		float64(0.2*c.gameState().ClientInfo.Resources)),
	)
	bestValue += shared.Resources(math.Min(
		rand.Float64(),
		float64(0.01*c.gameState().ClientInfo.Resources),
	))

	forageDecision := shared.ForageDecision{
		Type:         bestForageType,
		Contribution: bestValue,
	}
	c.Logf(
		"[Forage][Decision]: Decision %v | Expected ROI %v",
		forageDecision, bestROI)

	return forageDecision
}

func (c *client) DesperateForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Forage][Decision]: Desperate | Decision %v", forageDecision)
	return forageDecision
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.randomForageTurns {
		return c.RandomForage(), nil
	} else if c.gameState().ClientInfo.Resources <= c.config.desperationThreshold {
		return c.DesperateForage(), nil
	} else {
		return c.NormalForage(), nil
	}
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, revenue shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		contribution: forageDecision.Contribution,
		revenue:      revenue,
	})

	c.Logf(
		"[Forage][Update]: ForageType %v | Profit %v | Contribution %v | Actual ROI %v",
		forageDecision.Type,
		revenue-forageDecision.Contribution,
		forageDecision.Contribution,
		(revenue/forageDecision.Contribution)-1,
	)
}
