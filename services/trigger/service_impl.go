package trigger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rituraj-junglee/container-critic/models"
	"github.com/rituraj-junglee/container-critic/repo/reportconfig"
	reportsvc "github.com/rituraj-junglee/container-critic/services/report"
	"github.com/rituraj-junglee/container-critic/services/slacker"
	"github.com/slack-go/slack/socketmode"
)

type service struct {
	socketClient     *socketmode.Client
	reportsvc        reportsvc.Service
	slacksvc         slacker.Service
	reportconfigrepo reportconfig.Repository
}

func NewService(socketClient *socketmode.Client, reportsvc reportsvc.Service, slacksvc slacker.Service, reportconfigrepo reportconfig.Repository) Service {
	return &service{
		socketClient:     socketClient,
		reportsvc:        reportsvc,
		slacksvc:         slacksvc,
		reportconfigrepo: reportconfigrepo,
	}
}

func (s *service) TrivyTrigger(ctx context.Context, triggerReq models.TriggerRequest) (err error) {
	// Call the reportsvc to generate the report
	reportReq := models.ReportRequest{
		ClusterName: triggerReq.ClusterName,
		Namespace:   triggerReq.Namespace,
		Kind:        triggerReq.Kind,
		Findings:    triggerReq.Findings,
	}

	report, err := s.reportsvc.GenerateReport(ctx, reportReq)
	if err != nil {
		log.Println("Error generating report: ", err)
		return
	}

	// Send Messsage to the slack channel
	message := fmt.Sprintf("⚠️ Report for cluster:\n %s\n\nTotal Cost: %.2f\n\n ", reportReq.ClusterName, report.TotalCost)
	err = s.slacksvc.SendMessage(ctx, message)
	if err != nil {
		log.Println("Error sending message: ", err)
		return
	}

	// Send the report to the slack channel
	filename := fmt.Sprintf("report-%s.md", reportReq.ClusterName)
	err = s.slacksvc.SendFile(ctx, filename, report.Template)
	if err != nil {
		log.Println("Error sending file: ", err)
		return
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	// Save the report in the database
	reportConfig := models.ReportConfig{
		ReportID:     reportReq.ClusterName,
		Timestamp:    now,
		TargetConfig: report.TargetConfig,
	}
	err = s.reportconfigrepo.UpdateReportConfig(ctx, reportConfig)
	if err != nil {
		log.Println("Error saving report: ", err)
	}
	return

}
