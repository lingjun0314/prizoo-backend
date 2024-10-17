Prizoo 微服務系統
===

此為 prizoo 應用的後端系統，目前尚在開發中，整體採用微服務系統架構，搭配 Firebase 平台實現資料庫及登入系統
主要功能為實現抽獎活動整合，提供 api 給手機端 app 調用，使用 corn 執行定時任務。

涉及技術（現已使用及後續規劃）
---

Go/Gin/go-micro/protobuf/consul/Redis/RabbitMQ/ElasticSearch/firestore/firebase Authentication/Firebase storage/AWS EC2/docker

原本想使用 gRPC 來加快服務間的通信，但開發時 go-micro 框架剛更新至 v5 版本，還暫時沒有支援 gRPC ，因此選擇 http 搭配 protobuf 開發
