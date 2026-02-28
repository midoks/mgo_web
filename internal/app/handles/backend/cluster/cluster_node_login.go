package cluster

import (
	"time"

	"github.com/gin-gonic/gin"

	"mgo/internal/app/common"
	"mgo/internal/app/form"
	"mgo/internal/db"
	"mgo/internal/model"
)

func PostClusterNodeLoginAdd(c *gin.Context) {
	var field form.ClusterNodeLoginAdd
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	common_data := &model.ClusterNodeLogin{
		NodeID:     field.NodeID,
		UpdateTime: time.Now(),
	}

	if field.ID > 0 {
		if err := db.GetDb().Model(&model.ClusterNodeLogin{}).Where("id = ?", field.ID).Updates(common_data).Error; err != nil {
			common.ErrorResp(c, err, -1)
			return
		}
		common.SuccessResp(c)
		return
	}
	common_data.Status = true
	common_data.CreateTime = time.Now()
	if err := db.GetDb().Create(common_data).Error; err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
