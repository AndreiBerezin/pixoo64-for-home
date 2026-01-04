# pixoo64-for-home
Divoom pixoo64 display for home usage


<img src="static/readme/example.png" alt="example" width="448"/>


**run**
```
docker run -d \
  --name pixoo64 \
  --restart always \
  -e PIXOO_ADDRESS="192.168.0.100" \
  -e YANDEX_WEATHER_KEY="key" \
  -e YANDEX_WEATHER_LAT="55.751" \
  -e YANDEX_WEATHER_LON="37.618" \
  ghcr.io/andreiberezin/pixoo64:latest
```