package pkg

import (
	"fmt"
	"log"
	"time"

	rdApi "github.com/ophum/humstack-redeployment/pkg/api"
	rdClient "github.com/ophum/humstack-redeployment/pkg/client"
	"github.com/ophum/humstack/pkg/api/meta"
)

type CoopAgent struct {
	client         *ScoreServerClient
	redeployClient *rdClient.RedeploymentClient
	config         *Config
	lastSync       *time.Time
}

func NewCoopAgent(config *Config) (*CoopAgent, error) {
	client, err := NewScoreServerClient(config.ScoreServerURL)
	if err != nil {
		return nil, err
	}

	redeployClient := rdClient.NewRedeploymentClient("http", config.RedeploymentApiServerAddress, int32(config.RedeploymentApiServerPort))

	return &CoopAgent{
		client:         client,
		config:         config,
		lastSync:       nil,
		redeployClient: redeployClient,
	}, nil
}

func (a *CoopAgent) Run() {
	// first sync
	if err := a.sync(); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("sync done! `%s`", a.lastSync.Format(time.RFC3339Nano))
	}

	ticker := time.NewTicker(time.Second * time.Duration(a.config.PollingSeconds))
	defer ticker.Stop()

	// sync per PollingSeconds
	for {
		select {
		case <-ticker.C:
			if err := a.sync(); err != nil {
				log.Println(err.Error())
			} else {
				log.Printf("sync done! `%s`", a.lastSync.Format(time.RFC3339Nano))
			}
		}
	}
}
func (a *CoopAgent) sync() error {
	now := time.Now()
	if a.lastSync == nil {
		t, err := time.Parse(time.RFC3339Nano, "2021-02-01T00:00:00+09:00")
		if err != nil {
			return err
		}
		a.lastSync = &t
	}
	res, err := a.client.GetPenalties(*a.lastSync)
	if err != nil {
		// 認証エラー以外は終了する
		if err.Error() != "unauthorized" {
			return err
		}
		// 再認証
		if err := a.client.Authenticate(a.config.Name, a.config.Password); err != nil {
			return err
		}
		// 再取得
		res, err = a.client.GetPenalties(*a.lastSync)
		// ここでエラーの場合はどうしようもないので終了
		if err != nil {
			return err
		}
	}

	// score serverから取得したペナルティの一覧
	for _, p := range res.Data.Penalties {
		// すでにhumstack-redeploymentに登録しているか調べる
		rd, err := a.redeployClient.Get(p.ID)
		if err != nil {
			return err
		}
		// clientの仕様がｶｽなのでなくても空の場合notfound
		// 空でない場合はすでにsync済み
		if rd.ID != "" {
			continue
		}

		if _, err := a.redeployClient.Create(&rdApi.Redeployment{
			Meta: meta.Meta{
				ID: p.ID,
			},
			Spec: rdApi.RedeploymentSpec{
				Group:       a.config.Group,
				Namespace:   p.Problem.Code,
				VMIDPrefix:  fmt.Sprintf("team%02d_%s_", p.Team.Number, p.Problem.Code),
				RestartTime: p.CreatedAt.Add(time.Minute * time.Duration(a.config.PenaltyMinutes)),
			},
		}); err != nil {
			return err
		}
	}
	a.lastSync = &now
	return nil
}
