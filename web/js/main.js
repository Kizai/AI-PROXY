// 只保留API配置SPA逻辑
const pageRoutes = {
    'api-config': 'pages/api-config.html'
};

async function loadPage(pageName) {
    if (!checkAuth()) {
        showTokenModal();
        return;
    }
    try {
        showLoading();
        updateNavbar(pageName);
        const response = await fetch(pageRoutes[pageName]);
        if (!response.ok) throw new Error(`页面加载失败: ${response.status} - ${response.statusText}`);
        const html = await response.text();
        const mainContent = document.getElementById('main-content');
        if (mainContent) mainContent.innerHTML = html;
        initializePage(pageName);
    } catch (error) {
        const mainContent = document.getElementById('main-content');
        if (mainContent) mainContent.innerHTML = `<div class='alert alert-danger'>页面加载失败: ${error.message}</div>`;
    } finally {
        hideLoading();
    }
}

function updateNavbar(pageName) {
    document.querySelectorAll('#sidebar .components li').forEach(li => li.classList.remove('active'));
    const currentLink = document.querySelector(`#sidebar .components li a[onclick=\"loadPage('${pageName}')\"]`);
    if (currentLink) currentLink.parentElement.classList.add('active');
}

function initializePage(pageName) {
    if (pageName === 'api-config') {
        if (window.loadAPIConfigs) window.loadAPIConfigs();
    }
    // 首页API调用示例区块切换逻辑
    if (document.getElementById('vendor-select') && document.getElementById('lang-select')) {
      const vendorSelect = document.getElementById('vendor-select');
      const langSelect = document.getElementById('lang-select');
      const codeBlock = document.getElementById('api-demo-code');
      const codeTag = codeBlock.querySelector('code');
      const copyBtn = document.getElementById('copy-demo-btn');

      // 示例代码映射
      const demoCode = {
        openai: {
          python: `import openai\nopenai.api_base = \"https://aceproxy.xyz/openai/v1/chat/completions\"\nopenai.api_key = \"sk-你的key\"\nresponse = openai.ChatCompletion.create(model=\"gpt-3.5-turbo\", messages=[{\"role\": \"user\", \"content\": \"你好\"}])\nprint(response.choices[0].message.content)`,
          nodejs: `const openai = require('openai');\nopenai.apiBase = 'https://aceproxy.xyz/openai/v1/chat/completions';\nopenai.apiKey = 'sk-你的key';\nconst res = await openai.chat.completions.create({ model: 'gpt-3.5-turbo', messages: [{ role: 'user', content: '你好' }] });\nconsole.log(res.choices[0].message.content);`,
          curl: `curl -X POST https://aceproxy.xyz/openai/v1/chat/completions \\\n  -H \"Authorization: Bearer sk-你的key\" \\\n  -H \"Content-Type: application/json\" \\\n  -d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "你好"}]}'`
        },
        claude: {
          python: `import openai\nopenai.api_base = \"https://aceproxy.xyz/claude/v1/claude/generate\"\nopenai.api_key = \"azure-你的key\"\nresponse = openai.ChatCompletion.create(model=\"gpt-35-turbo\", messages=[{\"role\": \"user\", \"content\": \"你好\"}])\nprint(response.choices[0].message.content)`,
          nodejs: `const openai = require('openai');\nopenai.apiBase = 'https://aceproxy.xyz/claude/v1/claude/generate';\nopenai.apiKey = 'azure-你的key';\nconst res = await openai.chat.completions.create({ model: 'gpt-35-turbo', messages: [{ role: 'user', content: '你好' }] });\nconsole.log(res.choices[0].message.content);`,
          curl: `curl -X POST https://aceproxy.xyz/claude/v1/claude/generate \\\n  -H \"Authorization: Bearer azure-你的key\" \\\n  -H \"Content-Type: application/json\" \\\n  -d '{"model": "gpt-35-turbo", "messages": [{"role": "user", "content": "你好"}]}'`
        },
        grok: {
          python: `import requests\nurl = \"https://aceproxy.xyz/grok/v1/chat/completions\"\nheaders = {\"Authorization\": \"Bearer gr-你的key\", \"Content-Type\": \"application/json\"}\ndata = {\"model\": \"grok-1\", \"messages\": [{\"role\": \"user\", \"content\": \"你好\"}]}\nresponse = requests.post(url, headers=headers, json=data)\nprint(response.json())`,
          nodejs: `const axios = require('axios');\nconst res = await axios.post('https://aceproxy.xyz/grok/v1/chat/completions', {\n  model: 'grok-1',\n  messages: [{ role: 'user', content: '你好' }]\n}, {\n  headers: { Authorization: 'Bearer gr-你的key', 'Content-Type': 'application/json' }\n});\nconsole.log(res.data);`,
          curl: `curl -X POST https://aceproxy.xyz/grok/v1/chat/completions \\\n  -H \"Authorization: Bearer gr-你的key\" \\\n  -H \"Content-Type: application/json\" \\\n  -d '{"model": "grok-1", "messages": [{"role": "user", "content": "你好"}]}'`
        },
        gemini: {
          python: `import requests\nurl = \"https://aceproxy.xyz/gemini/v1beta/models/gemini-pro:generateContent\"\nheaders = {\"Authorization\": \"Bearer gm-你的key\", \"Content-Type\": \"application/json\"}\ndata = {\"model\": \"gemini-pro\", \"messages\": [{\"role\": \"user\", \"content\": \"你好\"}]}\nresponse = requests.post(url, headers=headers, json=data)\nprint(response.json())`,
          nodejs: `const axios = require('axios');\nconst res = await axios.post('https://aceproxy.xyz/gemini/v1beta/models/gemini-pro:generateContent', {\n  model: 'gemini-pro',\n  messages: [{ role: 'user', content: '你好' }]\n}, {\n  headers: { Authorization: 'Bearer gm-你的key', 'Content-Type': 'application/json' }\n});\nconsole.log(res.data);`,
          curl: `curl -X POST https://aceproxy.xyz/gemini/v1beta/models/gemini-pro:generateContent \\\n  -H \"Authorization: Bearer gm-你的key\" \\\n  -H \"Content-Type: application/json\" \\\n  -d '{"model": "gemini-pro", "messages": [{"role": "user", "content": "你好"}]}'`
        }
      };

      function updateDemoCode() {
        const vendor = vendorSelect.value;
        const lang = langSelect.value;
        const code = demoCode[vendor][lang] || '';
        codeTag.textContent = code;
        codeBlock.className = 'language-' + lang;
        if (window.Prism) Prism.highlightElement(codeTag);
      }

      vendorSelect.addEventListener('change', updateDemoCode);
      langSelect.addEventListener('change', updateDemoCode);

      // 复制按钮
      copyBtn.addEventListener('click', function() {
        navigator.clipboard.writeText(codeTag.textContent).then(function() {
          copyBtn.textContent = '已复制';
          setTimeout(() => { copyBtn.innerHTML = '<i class="bi bi-clipboard"></i>复制'; }, 1200);
        });
      });

      // 初始化
      updateDemoCode();
    }
}

// 认证、加载动画、错误提示等工具函数可保留在utils.js
// 其余无关功能全部移除
