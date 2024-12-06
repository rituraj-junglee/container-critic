package trigger

import (
	"context"
	"fmt"
	"log"

	"github.com/rituraj-junglee/container-critic/models"
	reportsvc "github.com/rituraj-junglee/container-critic/services/report"
	"github.com/rituraj-junglee/container-critic/services/slacker"
	"github.com/slack-go/slack/socketmode"
)

type service struct {
	socketClient *socketmode.Client
	reportsvc    reportsvc.Service
	slacksvc     slacker.Service
}

func NewService(socketClient *socketmode.Client, reportsvc reportsvc.Service, slacksvc slacker.Service) Service {
	return &service{
		socketClient: socketClient,
		reportsvc:    reportsvc,
		slacksvc:     slacksvc,
	}
}

func (s *service) TrivyTrigger(ctx context.Context, triggerReq models.TriggerRequest) (err error) {
	// Call the reportsvc to generate the report
	reportReq := models.ReportRequest{
		ClusterName: triggerReq.ClusterName,
		Findings:    triggerReq.Findings,
	}

	report, err := s.reportsvc.GenerateReport(ctx, reportReq)
	if err != nil {
		log.Println("Error generating report: ", err)
		return
	}

	// Send the report to the slack channel
	// s.slacksvc.SendMessage(ctx, report.Template)
	filename := fmt.Sprintf("report-%s.md", reportReq.ClusterName)
	err = s.slacksvc.SendFile(ctx, filename, report.Template)
	if err != nil {
		log.Println("Error sending file: ", err)
		return
	}

	// Send the report to the slack channel
	return

}
