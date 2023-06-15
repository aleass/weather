# 天气预报
通过彩云获取天气预报,并由企业微信发送通知

# 配置文件：yaml
wechat:
  - token: wechat token
    > wechat的url
  - note: xxxx
    > wechat的url的备注


caiyun:
  token : xxxxx
  > 彩云token
  addres:
  - name : xxxxx
  > 坐标名称
  - wechatNotes : xxxxx
  > 每个坐标发送的指定wechat note
  - coordinate : xxxxx,xxxxx
  > 坐标的经纬度
  - switch : bool
  > 是否启用

urlConfigPass:
  > url订阅配置
  - name: xxxx
  > 账号
    note: wechat.note
  > wechat的note