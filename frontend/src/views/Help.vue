<template>
  <div class="page-container">
    <div class="section-title">❓ 帮助</div>

    <div class="help-section">
      <h3>关于</h3>
      <div class="help-info">
        <p><strong>哔咔漫画 Web 版</strong> (PicACG Web)</p>
        <p>版本: 1.0.0</p>
        <p>基于 <a href="https://github.com/tonquer/picacg-qt" target="_blank">tonquer/picacg-qt</a> 的 Web 重构版本</p>
        <p>技术栈: Go + Vue 3 + TypeScript</p>
      </div>
    </div>

    <div class="help-section">
      <h3>功能说明</h3>
      <div class="help-info">
        <ul>
          <li>✅ 登录/注册哔咔漫画账号</li>
          <li>✅ 分类浏览漫画</li>
          <li>✅ 搜索漫画（支持热词）</li>
          <li>✅ 漫画详情查看</li>
          <li>✅ 在线漫画阅读器（单页/滚动模式）</li>
          <li>✅ 排行榜（24小时/本周/本月）</li>
          <li>✅ 收藏管理</li>
          <li>✅ 游戏区浏览</li>
          <li>✅ 聊天室</li>
          <li>✅ 锅贴广场</li>
          <li>✅ 下载管理</li>
          <li>✅ 个人资料/签到</li>
          <li>✅ 设置/修改密码</li>
          <li>✅ 阅读历史</li>
        </ul>
      </div>
    </div>

    <div class="help-section">
      <h3>使用说明</h3>
      <div class="help-info">
        <ol>
          <li>使用哔咔漫画账号登录（邮箱/用户名 + 密码）</li>
          <li>首次使用需要注册哔咔漫画账号（通过哔咔漫画 App 注册）</li>
          <li>登录后即可浏览、搜索、阅读漫画</li>
          <li>点击漫画封面查看详情，点击章节开始阅读</li>
          <li>阅读器支持单页模式和滚动模式切换</li>
          <li>可以将漫画添加到下载队列，下载到 NAS 存储</li>
        </ol>
      </div>
    </div>

    <div class="help-section">
      <h3>常见问题</h3>
      <div class="help-info">
        <div class="faq-item">
          <p class="faq-q">Q: 为什么登录失败？</p>
          <p class="faq-a">A: 请检查邮箱和密码是否正确。首次使用需要先在哔咔漫画 App 注册账号。</p>
        </div>
        <div class="faq-item">
          <p class="faq-q">Q: 为什么图片加载失败？</p>
          <p class="faq-a">A: 可能是网络问题或 PicACG 服务器暂时不可用，请稍后重试。</p>
        </div>
        <div class="faq-item">
          <p class="faq-q">Q: 下载的漫画保存在哪里？</p>
          <p class="faq-a">A: 下载的漫画保存在 NAS 的 /volume1/bika 目录下，按漫画名称和章节组织。</p>
        </div>
      </div>
    </div>

    <div class="help-section">
      <h3>相关链接</h3>
      <div class="help-info">
        <ul>
          <li><a href="https://github.com/Mogvl/bika" target="_blank">GitHub 仓库</a></li>
          <li><a href="https://github.com/tonquer/picacg-qt" target="_blank">原版 picacg-qt</a></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.help-section {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: var(--shadow);
}

.help-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--primary);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}

.help-info {
  font-size: 14px;
  line-height: 1.8;
  color: var(--text);
}

.help-info ul, .help-info ol {
  padding-left: 20px;
}

.help-info li {
  margin-bottom: 8px;
}

.help-info a {
  color: var(--primary);
  text-decoration: none;
}

.help-info a:hover {
  text-decoration: underline;
}

.faq-item {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}

.faq-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.faq-q {
  font-weight: 600;
  color: var(--primary);
  margin-bottom: 4px;
}

.faq-a {
  color: var(--text-secondary);
}
</style>
