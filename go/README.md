開啟後端 docker 的操作手冊
===

安裝 docker
---

首先需要確認 windows 下的 linux 子系統有沒有開啟，在開始搜尋 開啟或關閉windows功能

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E9%96%8B%E5%95%9F%E6%88%96%E9%97%9C%E9%96%89.png)

確認自己的這個選項有打開：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/linux%E5%AD%90%E7%B3%BB%E7%B5%B1.png)

同時要確保 CPU 虛擬化功能有打開，可以進入工作管理員查看：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E8%99%9B%E6%93%AC%E5%8C%96.png)

然後就可以去 docker 官網下載 docker desktop 了：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E4%B8%8B%E8%BC%89docker%20desktop.png)

打開 docker desktop 後登入共用帳號，就可以看到下面的畫面，點擊右上角的設定：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/docker%E8%A8%AD%E5%AE%9A.png)

在 General 標籤找到 "Use the WSL 2 based engine" 確保這個選項有勾選，如果沒有那需要勾選並重新啟動：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/docker%20general.png)

題外話，如果已經開啟過 docker desktop 然後關閉視窗，那麼他會在導覽列的小程式集裡面，再次點擊桌面捷徑或從開始菜單執行都是沒辦法打開的，只能從小程式集裡面打開：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E5%B0%8F%E7%A8%8B%E5%BC%8F%E9%9B%86.png)

拉取鏡像
---

在這個專案中一共會用到兩個（未來會有三個）公共鏡像，分別是 redis 及 hashicorp/consul，從最上方的搜尋欄輸入這兩個公共鏡像的名字（圖片以 Redis 舉例）：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E6%8B%89%E5%AE%98%E6%96%B9%E9%8F%A1%E5%83%8F.png)

圖片中第一個有小綠標的鏡像就是要用的了，點擊 pull 進行拉取，hashicorp/consul 也是同樣的道理
這時候切換到 images 標籤就可以看到拉取的兩個鏡像了：

![test](https://github.com/lingjun0314/prizoo-system/blob/main/go/images/%E6%8B%89%E5%8F%96%E7%9A%84%E9%8F%A1%E5%83%8F.png)

保留在這個頁面，將標籤從 local 切換成 hub，列表中就是所有後端的服務鏡像，如果看到裡面有很多個鏡像，那麼就會需要 pull 所有 tag 為 latest 的鏡像
（待更新）
