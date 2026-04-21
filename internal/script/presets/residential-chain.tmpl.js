// 预设脚本：链式代理（起飞节点 → 住宅 IP 落地）
//
// 作用：
//   - 把机场订阅节点聚合成「🛫 AI起飞节点」select 组（供链式代理第一跳）
//   - 把用户填写的住宅 IP 代理聚合成「🛬 AI落地节点」select 组
//     （底层用 mihomo 的 dialer-proxy 实现 AI起飞 → 住宅IP 链路）
//   - 把 AI/开发工具相关域名（Claude/ChatGPT/Cursor/Termius 等）路由到 AI落地节点
//   - 阿里系域名直连，AI 服务始终走住宅 IP 不踩风控
//
// 住宅 IP 信息由 lan-proxy-gateway 渲染时注入。
// 其它字段（CUSTOM_INBOUND_PROXIES / INJECT_INTO_GROUPS）留空，
// 需要手改可取消下面相应注释。

function main(config) {

  // =============================================
  // 用户配置区
  // =============================================

  // 自定义起飞节点：默认空数组。
  // 若想把额外的 HTTP/SOCKS5 节点注入到机场的代理组里，取消下面注释并填写：
  //
  // const CUSTOM_INBOUND_PROXIES = [
  //   {
  //     name: "🌐 示例-HTTP",
  //     type: "http",
  //     server: "proxy.example.com",
  //     port: 8080,
  //     username: "user",
  //     password: "pass",
  //     tls: true,
  //     "skip-cert-verify": true
  //   },
  // ];
  const CUSTOM_INBOUND_PROXIES = [];

  // 住宅 IP 落地节点：由 lan-proxy-gateway 的 S → 链式代理预设 填好注入。
  // 这里不要手改 —— 通过 gateway 主菜单 → 3 换代理源 → S 增强脚本 重填。
  const RESIDENTIAL_PROXIES = [
    __RESIDENTIAL_PROXY_JSON__
  ];

  // 把自定义节点注入到哪些机场代理组（需要精确匹配组名）。默认空。
  // 取消注释例：
  //
  // const INJECT_INTO_GROUPS = [
  //   "🚀 节点选择",
  //   "♻️ 自动选择",
  // ];
  const INJECT_INTO_GROUPS = [];

  // =============================================
  // 工具函数
  // =============================================

  function ensureArray(value) {
    return Array.isArray(value) ? value : [];
  }

  function removeNamedItems(list, names) {
    const nameSet = new Set(names);
    return list.filter(item => !(item && item.name && nameSet.has(item.name)));
  }

  function uniqueStrings(list) {
    return Array.from(new Set(list.filter(Boolean)));
  }

  function mergeRules(priorityRules, existingRules) {
    const result = [];
    const seen = new Set();
    priorityRules.concat(existingRules).forEach(rule => {
      if (typeof rule !== "string" || seen.has(rule)) return;
      seen.add(rule);
      result.push(rule);
    });
    return result;
  }

  // =============================================
  // 主体逻辑
  // =============================================

  config.proxies = ensureArray(config.proxies);
  config["proxy-groups"] = ensureArray(config["proxy-groups"]);
  config.rules = ensureArray(config.rules);

  const customInboundNames = CUSTOM_INBOUND_PROXIES.map(p => p.name);
  const residentialNames = RESIDENTIAL_PROXIES.map(p => p.name);
  const allManagedNames = [
    ...customInboundNames,
    ...residentialNames,
    "🛫 AI起飞节点",
    "🛬 AI落地节点"
  ];

  // 清理旧的同名节点，避免重复注入
  config.proxies = removeNamedItems(config.proxies, [...customInboundNames, ...residentialNames]);

  // 注入节点到 config.proxies
  CUSTOM_INBOUND_PROXIES.forEach(p => config.proxies.push(p));
  RESIDENTIAL_PROXIES.forEach(p => config.proxies.push(p));

  // 收集机场订阅里原有的节点名（排除自己管理的）
  const myProxyNames = new Set([...customInboundNames, ...residentialNames]);
  const subscriptionProxyNames = uniqueStrings(
    config.proxies
      .map(p => p && p.name)
      .filter(name => name && !myProxyNames.has(name))
  );

  // 🛫 AI起飞节点 = 订阅节点 + 自定义起飞节点
  const inboundProxies = uniqueStrings([
    ...subscriptionProxyNames,
    ...customInboundNames
  ]);

  // 🛬 AI落地节点 = 住宅 IP 节点 + 起飞节点（链式代理兜底）
  const outboundProxies = [
    ...residentialNames,
    "🛫 AI起飞节点",
  ];

  // 清理旧分组，避免重复注入
  config["proxy-groups"] = removeNamedItems(config["proxy-groups"], allManagedNames);

  config["proxy-groups"].unshift(
    {
      name: "🛫 AI起飞节点",
      type: "select",
      proxies: inboundProxies.length ? inboundProxies : ["DIRECT"]
    },
    {
      name: "🛬 AI落地节点",
      type: "select",
      proxies: outboundProxies
    }
  );

  // 把自定义节点注入到指定的机场代理组
  if (INJECT_INTO_GROUPS.length > 0) {
    config["proxy-groups"].forEach(group => {
      if (INJECT_INTO_GROUPS.includes(group.name)) {
        group.proxies = uniqueStrings([
          ...customInboundNames,
          ...ensureArray(group.proxies)
        ]);
      }
    });
  }

  // =============================================
  // 优先规则
  // =============================================

  const priorityRules = [
    // 阿里系 → 直连
    "DOMAIN-SUFFIX,alibaba-inc.com,DIRECT",
    "DOMAIN-SUFFIX,alibaba.net,DIRECT",
    "DOMAIN-SUFFIX,dingtalk.com,DIRECT",
    "DOMAIN-SUFFIX,alipay.com,DIRECT",
    "DOMAIN-SUFFIX,aliyun.com,DIRECT",
    "PROCESS-NAME,DingTalk,DIRECT",
    "PROCESS-NAME,iDingTalk,DIRECT",

    // Claude / Anthropic
    "DOMAIN-SUFFIX,anthropic.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,claude.ai,🛬 AI落地节点",
    "DOMAIN-SUFFIX,claude.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,claudeusercontent.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,usefathom.com,🛬 AI落地节点",
    "PROCESS-NAME,Claude,🛬 AI落地节点",

    // Cursor
    "DOMAIN-SUFFIX,cursor.sh,🛬 AI落地节点",
    "DOMAIN-SUFFIX,cursor.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,cursor-cdn.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,cursorapi.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,cursor.so,🛬 AI落地节点",
    "DOMAIN,api2.cursor.sh,🛬 AI落地节点",
    "DOMAIN,api3.cursor.sh,🛬 AI落地节点",
    "DOMAIN,repo.cursor.sh,🛬 AI落地节点",
    "DOMAIN,download.cursor.sh,🛬 AI落地节点",
    "DOMAIN,marketplace.cursor.sh,🛬 AI落地节点",
    "PROCESS-NAME,Cursor,🛬 AI落地节点",
    "PROCESS-NAME,Cursor Helper,🛬 AI落地节点",
    "PROCESS-NAME,Cursor Helper (GPU),🛬 AI落地节点",
    "PROCESS-NAME,Cursor Helper (Plugin),🛬 AI落地节点",
    "PROCESS-NAME,Cursor Helper (Renderer),🛬 AI落地节点",

    // OpenAI / ChatGPT
    "DOMAIN-SUFFIX,openai.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,chatgpt.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,oaistatic.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,oaiusercontent.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,openaiapi.com,🛬 AI落地节点",
    "DOMAIN,cdn.openai.com,🛬 AI落地节点",
    "PROCESS-NAME,ChatGPT,🛬 AI落地节点",

    // 通用 AI 基础设施
    "DOMAIN-SUFFIX,auth0.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,intercom.io,🛬 AI落地节点",
    "DOMAIN-SUFFIX,intercomcdn.com,🛬 AI落地节点",

    // IP 检测
    "DOMAIN-SUFFIX,ping0.cc,🛬 AI落地节点",
    "DOMAIN,checkip.amazonaws.com,🛬 AI落地节点",
    "DOMAIN,ip.sb,🛬 AI落地节点",
    "DOMAIN,ipinfo.io,🛬 AI落地节点",
    "DOMAIN,ipapi.co,🛬 AI落地节点",

    // Termius
    "DOMAIN,api.termius.com,🛬 AI落地节点",
    "DOMAIN,grpc.termius.com,🛬 AI落地节点",
    "DOMAIN,account.termius.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,termius.com,🛬 AI落地节点",
    "DOMAIN-SUFFIX,termi.us,🛬 AI落地节点",
    "PROCESS-NAME,Termius,🛬 AI落地节点",
    "PROCESS-NAME,Termius Helper,🛬 AI落地节点",

    // 其他 AI 产品（按需取消注释）
    // "DOMAIN-SUFFIX,gemini.google.com,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,perplexity.ai,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,midjourney.com,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,replicate.com,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,huggingface.co,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,cohere.ai,🛬 AI落地节点",
    // "DOMAIN-SUFFIX,groq.com,🛬 AI落地节点",
  ];

  config.rules = mergeRules(priorityRules, config.rules);

  return config;
}
