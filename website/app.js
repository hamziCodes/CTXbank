// VERTEX/ctxbank site logic: tabs, copy buttons, theme, Connect page.
(function () {
  'use strict';

  var VALID_TABS = ['how', 'setup', 'connect'];

  /* ---------- tabs (hash deep-linkable: #how #setup #connect) ---------- */
  var tabButtons = Array.prototype.slice.call(document.querySelectorAll('[data-tab]'));

  function activateTab(name, pushHash) {
    if (VALID_TABS.indexOf(name) === -1) name = 'how';
    tabButtons.forEach(function (btn) {
      var active = btn.getAttribute('data-tab') === name;
      btn.classList.toggle('active', active);
      btn.setAttribute('aria-selected', active ? 'true' : 'false');
    });
    VALID_TABS.forEach(function (t) {
      var panel = document.getElementById('panel-' + t);
      if (!panel) return;
      var show = t === name;
      panel.classList.toggle('active', show);
      if (show) {
        panel.removeAttribute('hidden');
        // re-trigger entrance animation
        panel.style.animation = 'none';
        void panel.offsetWidth;
        panel.style.animation = '';
      } else {
        panel.setAttribute('hidden', '');
      }
    });
    if (pushHash !== false) {
      history.replaceState(null, '', '#' + name);
    }
    window.scrollTo({ top: 0, behavior: 'smooth' });
    observeReveals();
  }

  tabButtons.forEach(function (btn) {
    btn.addEventListener('click', function () { activateTab(btn.getAttribute('data-tab')); });
  });
  document.querySelectorAll('[data-tab-link]').forEach(function (el) {
    el.addEventListener('click', function (e) {
      e.preventDefault();
      activateTab(el.getAttribute('data-tab-link'));
    });
  });
  window.addEventListener('hashchange', function () {
    activateTab(location.hash.replace('#', ''), false);
  });

  /* ---------- copy buttons ---------- */
  document.querySelectorAll('.copy-btn').forEach(function (btn) {
    btn.addEventListener('click', function () {
      var text = btn.getAttribute('data-copy');
      function done() {
        btn.classList.add('copied');
        var orig = btn.textContent;
        btn.textContent = 'Copied';
        setTimeout(function () { btn.classList.remove('copied'); btn.textContent = orig; }, 1600);
      }
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(done, done);
      } else {
        var ta = document.createElement('textarea');
        ta.value = text;
        document.body.appendChild(ta);
        ta.select();
        try { document.execCommand('copy'); } catch (e) {}
        document.body.removeChild(ta);
        done();
      }
    });
  });

  /* ---------- theme ---------- */
  var themeBtn = document.getElementById('theme-toggle');
  function applyTheme(t) {
    document.documentElement.setAttribute('data-theme', t);
    try { localStorage.setItem('ctxbank-theme', t); } catch (e) {}
  }
  var savedTheme = null;
  try { savedTheme = localStorage.getItem('ctxbank-theme'); } catch (e) {}
  applyTheme(savedTheme || 'light');
  if (themeBtn) {
    themeBtn.addEventListener('click', function () {
      var cur = document.documentElement.getAttribute('data-theme');
      applyTheme(cur === 'dark' ? 'light' : 'dark');
    });
  }

  /* ---------- reveal on scroll ---------- */
  var observer = null;
  function observeReveals() {
    if (!('IntersectionObserver' in window)) {
      document.querySelectorAll('.reveal').forEach(function (el) { el.classList.add('visible'); });
      return;
    }
    if (!observer) {
      observer = new IntersectionObserver(function (entries) {
        entries.forEach(function (en) {
          if (en.isIntersecting) { en.target.classList.add('visible'); observer.unobserve(en.target); }
        });
      }, { threshold: 0.12 });
    }
    document.querySelectorAll('.tab-panel.active .reveal:not(.visible)').forEach(function (el) { observer.observe(el); });
  }

  /* ---------- Connect page ---------- */
  var tokenInput = document.getElementById('connect-token');
  var portInput = document.getElementById('connect-port');
  var testBtn = document.getElementById('btn-test-connection');
  var openBtn = document.getElementById('btn-open-dashboard');
  var statusBox = document.getElementById('connect-status');
  var lastGoodUrl = null;

  function localUrl() {
    var port = parseInt(portInput.value, 10) || 4242;
    var token = tokenInput.value.trim();
    return 'http://localhost:' + port + '?token=' + encodeURIComponent(token);
  }

  function setStatus(kind, html) {
    statusBox.className = 'connect-status' + (kind ? ' ' + kind : '');
    statusBox.innerHTML = html;
  }

  // Returns a promise resolving true if SOMETHING answers on the port.
  // A no-cors fetch resolves (opaque) even when the server sends no CORS
  // headers — which is exactly what an outdated ctx does. That lets us
  // tell "old version" apart from "nothing running".
  function probeReachable(port, ms) {
    return new Promise(function (resolve) {
      var controller = new AbortController();
      var timer = setTimeout(function () { controller.abort(); resolve(false); }, ms || 2500);
      fetch('http://localhost:' + port + '/api/status', { mode: 'no-cors', signal: controller.signal })
        .then(function () { clearTimeout(timer); resolve(true); })
        .catch(function () { clearTimeout(timer); resolve(false); });
    });
  }

  function testConnection() {
    var token = tokenInput.value.trim();
    if (!token) {
      setStatus('err', '<span class="status-dot err"></span>Paste your project token first — run <code>ctx token</code> in your project to get it.');
      return;
    }
    var port = parseInt(portInput.value, 10) || 4242;
    var url = 'http://localhost:' + port + '/api/status?token=' + encodeURIComponent(token);
    setStatus('', '<span class="status-dot busy"></span>Knocking on your local dashboard…');
    openBtn.disabled = true;
    lastGoodUrl = null;

    var controller = new AbortController();
    var timer = setTimeout(function () { controller.abort(); }, 5000);

    fetch(url, { signal: controller.signal, mode: 'cors' })
      .then(function (res) {
        clearTimeout(timer);
        if (res.status === 401) throw new Error('bad-token');
        if (!res.ok) throw new Error('http-' + res.status);
        return res.json();
      })
      .then(function (data) {
        var name = (data && data.project_name) || 'your project';
        var branch = data && data.branch ? ' &middot; <code>' + escapeHtml(data.branch) + '</code>' : '';
        var ver = data && data.ctx_version ? ' &middot; ctx v' + escapeHtml(data.ctx_version) : '';
        lastGoodUrl = localUrl();
        openBtn.disabled = false;
        setStatus('ok', '<span class="status-dot ok"></span>Connected to <strong>' + escapeHtml(name) + '</strong>' + branch + ver + '. Your dashboard is ready.');
      })
      .catch(function (err) {
        clearTimeout(timer);
        if (err && err.message === 'bad-token') {
          setStatus('err', '<span class="status-dot err"></span>That token was rejected. Make sure it matches <code>ctx token</code> output <em>in the same project folder</em> where you ran <code>ctx ui</code> — every project has its own token.');
          return;
        }
        // Network-level failure: is anything even listening?
        setStatus('', '<span class="status-dot busy"></span>That didn\'t answer — checking what\'s on that port…');
        probeReachable(port).then(function (reachable) {
          if (reachable) {
            setStatus('err', '<span class="status-dot err"></span>Something is running on port ' + port + ', but it won\'t talk to this page. You\'re almost certainly on an <strong>outdated ctx</strong> (older versions can\'t do this handshake). Update: <code>go install github.com/hamziCodes/CTXbank/cmd/ctx@latest</code> — or re-run the installer — then <code>ctx ui</code> again.');
          } else {
            setStatus('err', '<span class="status-dot err"></span>Can\'t reach your dashboard. Is <code>ctx ui</code> running in your project? It serves on <code>localhost:' + port + '</code>.');
          }
        });
      });
  }

  function escapeHtml(s) {
    return String(s).replace(/[&<>"']/g, function (c) {
      return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c];
    });
  }

  if (testBtn) testBtn.addEventListener('click', testConnection);
  if (openBtn) openBtn.addEventListener('click', function () {
    if (lastGoodUrl) window.open(lastGoodUrl, '_blank', 'noopener');
  });
  if (tokenInput) {
    // Typing a new token invalidates the previous check.
    tokenInput.addEventListener('input', function () {
      openBtn.disabled = true;
      lastGoodUrl = null;
    });
    tokenInput.addEventListener('keydown', function (e) {
      if (e.key === 'Enter') testConnection();
    });
  }

  /* ---------- init ---------- */
  activateTab(location.hash.replace('#', '') || 'how', false);
})();
