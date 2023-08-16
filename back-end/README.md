# Usage

Để chạy project, ta thực hiện lần lượt các lệnh:

```
docker-compose build

docker-compose up 
```


# API Document

### Access Path

```
http://localhost:8080
```


### API Endpoints

| HTTP Verbs | Endpoints | Action |
| --- | --- | --- |
| GET | /api/rate/latest | Lấy rate ngày gần nhất |
| GET | /api/rate/{date} | Lấy rate của một ngày(YYYY-MM-DD) nào đó |
| GET | /api/rate/analyze | Lấy các giá trị thấp nhất, cáo nhất, trung bình của từng loại đồng |



### Functions 

- GetRatesLatest() => Lấy rate ngày mới nhất

- GetRatesByDate() => Lấy rate của một ngày(YYYY-MM-DD) nào đó

- GetRatesAnalyze() => Lấy các giá trị thấp nhất, cáo nhất, trung bình của từng loại đồng