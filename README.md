# 豆瓣书影音 App 前端

这是一个仿照豆瓣的书影音管理应用前端，使用 Vue 3 + TypeScript + Element Plus 构建。

## 功能特性

- 🎵 **音乐管理**：添加、编辑、删除专辑信息
- 📚 **书籍管理**：记录书籍信息，包括作者、类型、评分等
- 🎬 **影视管理**：管理电影和电视剧信息
- 🔍 **CRUD 操作**：完整的增删改查功能
- 🎨 **响应式设计**：适配桌面和移动设备
- ⭐ **评分系统**：0-5 星评分
- 🖼️ **图片支持**：支持封面图片显示

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript
- **Element Plus** - Vue 3 组件库
- **Vite** - 构建工具

## 项目结构

```
src/
├── components/          # Vue 组件
│   ├── MainPage.vue    # 主页面（导航）
│   ├── Music.vue       # 音乐页面
│   ├── Books.vue       # 书籍页面
│   └── Movies.vue      # 影视页面
├── services/           # API 服务
│   ├── api.ts         # API 接口
│   └── types.ts       # TypeScript 类型定义
├── App.vue            # 根组件
└── main.ts            # 入口文件
```

## 安装和运行

### 前置要求

确保后端服务正在运行（默认端口 8080）

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
pnpm dev
```

应用将在 `http://localhost:5173` 启动

### 构建生产版本

```bash
pnpm build
```

## API 配置

前端默认连接到 `http://localhost:8080`，如果需要修改后端地址，请编辑：

```typescript
// src/services/api.ts
const API_BASE_URL = 'http://your-backend-url:port'
```

## 使用说明

### 添加条目
1. 点击页面右上角的 "+" 按钮
2. 填写必填字段（标题和艺术家/作者/导演）
3. 可选填写其他信息（类型、年份、评分、链接、备注等）
4. 对于音乐，可以添加曲目列表

### 编辑条目
1. 点击条目卡片上的"编辑"按钮
2. 修改需要更新的字段
3. 点击"更新"保存更改

### 删除条目
1. 点击条目卡片上的"删除"按钮
2. 确认删除操作

## 数据模型

### 音乐 (Album)
```typescript
interface Album {
  title: string      // 标题（必需）
  artist: string     // 艺术家（必需）
  genre: string      // 流派
  year: number       // 年份
  cuts: string[]     // 曲目列表
  url: string        // 链接
  artwork: string    // 封面图片链接
  comment: string    // 备注
  rating: number     // 评分 (0-5)
}
```

### 书籍 (Book)
```typescript
interface Book {
  title: string      // 标题（必需）
  author: string     // 作者（必需）
  genre: string      // 类型
  year: number       // 年份
  url: string        // 链接
  cover: string      // 封面图片链接
  comment: string    // 备注
  rating: number     // 评分 (0-5)
}
```

### 影视 (Movie)
```typescript
interface Movie {
  title: string      // 标题（必需）
  director: string   // 导演（必需）
  genre: string      // 类型
  year: number       // 年份
  url: string        // 链接
  comment: string    // 备注
  rating: number     // 评分 (0-5)
}
```

## 注意事项

1. **索引字段不可修改**：标题和艺术家/作者/导演字段在编辑时不可修改，因为这些字段用作唯一标识
2. **图片显示**：支持外部图片链接，如果图片加载失败会显示占位图标
3. **响应式设计**：在移动设备上会自动调整布局
4. **错误处理**：包含完整的错误提示和加载状态

## 开发说明

项目使用 Vue 3 的 Composition API 和 `<script setup>` 语法，结合 TypeScript 提供类型安全。

主要特性：
- 组件化开发
- 响应式状态管理
- TypeScript 类型定义
- Element Plus UI 组件
- 异步数据加载
- 错误边界处理

## 许可证

MIT License
