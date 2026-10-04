document.addEventListener('DOMContentLoaded', async () => {
  const status = document.getElementById('status');
  const toggle = document.getElementById('toggle');

  chrome.runtime.sendMessage({ type: 'getStatus' }, (response) => {
    if (!response) return;
    status.textContent = response.enabled ? 'Guard is active' : 'Guard is paused';
    toggle.textContent = response.enabled ? 'Pause guard' : 'Resume guard';
  });

  toggle.addEventListener('click', () => {
    chrome.runtime.sendMessage({ type: 'getStatus' }, (response) => {
      const next = !response.enabled;
      chrome.runtime.sendMessage({ type: 'setStatus', enabled: next }, () => {
        status.textContent = next ? 'Guard is active' : 'Guard is paused';
        toggle.textContent = next ? 'Pause guard' : 'Resume guard';
      });
    });
  });
});
