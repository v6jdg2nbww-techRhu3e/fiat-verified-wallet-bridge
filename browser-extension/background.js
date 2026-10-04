const BLOCKLIST = [
  "example-malicious-site.com",
  "example-phishing-page.net",
  "example-fake-login.org"
];

async function sendEvent(name, details) {
  try {
    await fetch('http://localhost:8080/api/monitor/submit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ event: name, details, source: 'owner-guard' })
    });
  } catch (err) {
    console.error('owner guard submit failed', err);
  }
}

chrome.webNavigation.onCommitted.addListener(async (details) => {
  const url = details.url || '';
  for (const bad of BLOCKLIST) {
    if (url.includes(bad)) {
      await sendEvent('blocked_domain', { url, reason: 'blocked by local list' });
      chrome.tabs.update(details.tabId, { url: 'about:blank' });
      return;
    }
  }
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'getStatus') {
    chrome.storage.local.get(['guardEnabled'], (result) => {
      sendResponse({ enabled: result.guardEnabled !== false });
    });
    return true;
  }

  if (message.type === 'setStatus') {
    chrome.storage.local.set({ guardEnabled: !!message.enabled }, () => {
      sendResponse({ success: true, enabled: !!message.enabled });
    });
    return true;
  }

  sendResponse({ ok: true });
  return true;
});
