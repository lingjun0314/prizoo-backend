Prizoo 微服務系統
===

此為 prizoo 應用的後端系統，目前尚在開發中，整體採用微服務系統架構，搭配 Firebase 平台實現資料庫及登入系統
主要功能為實現抽獎活動整合，提供 api 給手機端 app 調用，使用 corn 執行定時任務。

涉及技術（現已使用及後續規劃）
---

Go/Gin/go-micro/protobuf/consul/Redis/RabbitMQ/ElasticSearch/firestore/firebase Authentication/Firebase storage/AWS EC2/docker

原本想使用 gRPC 來加快服務間的通信，但開發時 go-micro 框架剛更新至 v5 版本，還暫時沒有支援 gRPC ，因此選擇 http 搭配 protobuf 開發

系統架構
---

整體體統資料流可以參考下圖：
![data flow](https://github.com/lingjun0314/prizoo-backend/blob/main/go/images/%E6%8A%BD%E7%8D%8E%E5%B9%B3%E5%8F%B0system%20structure.png)

目前微服務架構暫時拆分為七個服務，主要拆分依據為系統功能，且盡量避免不同服務操作相同資料表，具體架構如下：
![data flow](https://github.com/lingjun0314/prizoo-backend/blob/main/go/images/Microservices%20structure.drawio.png)
