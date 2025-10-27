# 验证码配置说明

本文档说明如何在 `core.yaml` 中配置验证码功能。

## 预置字体文件

```
// fonts/3Dumb.ttf (142.224kB)
// fonts/ApothecaryFont.ttf (62.08kB)
// fonts/Comismsh.ttf (80.132kB)
// fonts/DENNEthree-dee.ttf (83.188kB)
// fonts/DeborahFancyDress.ttf (32.52kB)
// fonts/Flim-Flam.ttf (140.576kB)
// fonts/RitaSmith.ttf (31.24kB)
// fonts/actionj.ttf (34.944kB)
// fonts/chromohv.ttf (45.9kB)
// fonts/wqy-microhei.ttc (5.177MB)
```

## 配置结构

```yaml
# 验证码配置
CaptchaConf:
  # 验证码类型: digit(数字), string(字符串), math(数学), chinese(中文), audio(音频), random(随机)
  Type: random

  # 数字验证码配置
  Digit:
    Height: 80 # 图片高度
    Width: 240 # 图片宽度
    Length: 5 # 验证码长度
    MaxSkew: 0.7 # 最大倾斜度
    DotCount: 80 # 干扰点数量

  # 字符串验证码配置
  String:
    Height: 60
    Width: 240
    NoiseCount: 0
    ShowLineOptions: 0
    Length: 4
    Source: "1234567890qwertyuioplkjhgfdsazxcvbnm" # 字符源
    BgColor: # 背景颜色 (RGBA)
      R: 254
      G: 254
      B: 254
      A: 254
    Fonts: # 字体列表
      - "wqy-microhei.ttc"
      - "3Dumb.ttf"
      - "ApothecaryFont.ttf"
      - "Comismsh.ttf"
      - "DENNEthree-dee.ttf"
      - "DeborahFancyDress.ttf"
      - "Flim-Flam.ttf"
      - "RitaSmith.ttf"
      - "actionj.ttf"
      - "chromohv.ttf"

  # 数学验证码配置
  Math:
    Height: 60
    Width: 240
    NoiseCount: 0
    ShowLineOptions: 0
    BgColor:
      R: 254
      G: 254
      B: 254
      A: 254
    Fonts:
      - "wqy-microhei.ttc"
      - "3Dumb.ttf"
      - "ApothecaryFont.ttf"
      - "Comismsh.ttf"
      - "DENNEthree-dee.ttf"
      - "DeborahFancyDress.ttf"
      - "Flim-Flam.ttf"
      - "RitaSmith.ttf"
      - "actionj.ttf"
      - "chromohv.ttf"

  # 中文验证码配置
  Chinese:
    Height: 60
    Width: 240
    NoiseCount: 0
    ShowLineOptions: 0
    Length: 2
    Source: "设想你在处理消费者的音频输出..." # 中文字符源
    BgColor:
      R: 254
      G: 254
      B: 254
      A: 254
    Fonts:
      - "wqy-microhei.ttc"
      - "3Dumb.ttf"
      - "ApothecaryFont.ttf"
      - "Comismsh.ttf"
      - "DENNEthree-dee.ttf"
      - "DeborahFancyDress.ttf"
      - "Flim-Flam.ttf"
      - "RitaSmith.ttf"
      - "actionj.ttf"
      - "chromohv.ttf"

  # 音频验证码配置
  Audio:
    Length: 4 # 音频验证码长度
    Language: "zh" # 语言：zh(中文), en(英文)

  # 随机验证码配置
  Random:
    # 启用的验证码类型
    EnabledTypes:
      - "digit"
      - "string"
      - "math"
      - "chinese"
    # 是否排除音频验证码
    ExcludeAudio: true

  # 存储配置
  Store:
    # 存储类型: memory(内存), redis(Redis)
    Type: redis
    # 过期时间 (分钟)
    Expire: 5
    # Key前缀
    KeyPrefix: "captcha:"
    # Redis配置
    Redis:
      Addr: "localhost:6379"
      Password: ""
      DB: 0
      PoolSize: 10
```

## 配置说明

### 验证码类型 (Type)

- `digit`: 数字验证码 (如: 12345)
- `string`: 字符串验证码 (如: abc123)
- `math`: 数学验证码 (如: 3+5=?)
- `chinese`: 中文验证码 (如: 你好)
- `audio`: 音频验证码
- `random`: 随机类型验证码

### 随机验证码配置 (Random)

当 `Type` 设置为 `random` 时，系统会从 `EnabledTypes` 中随机选择一种类型生成验证码。

- `EnabledTypes`: 启用的验证码类型列表
- `ExcludeAudio`: 是否排除音频验证码（推荐设置为 true，因为音频验证码可能不适合所有场景）

### 存储配置 (Store)

- `memory`: 内存存储，适合开发和测试环境
- `redis`: Redis 存储，适合生产环境

### 使用示例

在代码中使用验证码服务：

```go
// 在 ServiceContext 中已经初始化了 CaptchaService
func (l *SomeLogic) GenerateCaptcha() (*types.CaptchaResp, error) {
    // 生成验证码
    result, err := l.svcCtx.CaptchaService.Generate()
    if err != nil {
        return nil, err
    }

    return &types.CaptchaResp{
        ID:         result.ID,
        Base64Blob: result.Base64Blob,
    }, nil
}

// 验证验证码
func (l *SomeLogic) VerifyCaptcha(id, answer string) bool {
    return l.svcCtx.CaptchaService.VerifyAndClear(id, answer)
}
```

### 环境配置建议

**开发环境:**

```yaml
CaptchaConf:
  Type: digit
  Store:
    Type: memory
    Expire: 10
```

**生产环境:**

```yaml
CaptchaConf:
  Type: random
  Random:
    EnabledTypes: ["digit", "string", "math"]
    ExcludeAudio: true
  Store:
    Type: redis
    Expire: 5
    Redis:
      Addr: "your-redis-host:6379"
      Password: "your-redis-password"
```

## 注意事项

1. 字体文件需要确保在系统中存在，或者使用 base64Captcha 的默认字体
2. Redis 配置需要确保 Redis 服务可用
3. 音频验证码会生成较大的 Base64 数据，建议在带宽有限的环境中排除
4. 随机验证码增加了安全性，推荐在生产环境中使用
