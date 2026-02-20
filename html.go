package main

var indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
<title>CronPanel</title>
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Crect width='32' height='32' rx='8' fill='%234f46e5'/%3E%3Ccircle cx='16' cy='16' r='9' fill='none' stroke='white' stroke-width='2'/%3E%3Cline x1='16' y1='16' x2='16' y2='10' stroke='white' stroke-width='2.2' stroke-linecap='round'/%3E%3Cline x1='16' y1='16' x2='20' y2='18' stroke='%23a5b4fc' stroke-width='2.2' stroke-linecap='round'/%3E%3Ccircle cx='16' cy='16' r='2' fill='white'/%3E%3C/svg%3E">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap" rel="stylesheet">
<style>
:root{
  --bg:#f1f4f9;--sidebar-bg:#1e1b4b;--sidebar-border:rgba(255,255,255,0.08);
  --surface:#fff;--surface2:#f8fafc;--border:#e2e8f0;--border2:#cbd5e1;
  --accent:#4f46e5;--accent2:#7c3aed;--accent-light:#eef2ff;--accent-text:#4338ca;
  --green:#059669;--red:#dc2626;--red-bg:#fef2f2;
  --text:#0f172a;--text2:#475569;--text3:#94a3b8;--sidebar-text:rgba(255,255,255,0.65);
  --shadow:0 1px 3px rgba(0,0,0,0.06),0 1px 2px rgba(0,0,0,0.04);
  --shadow-md:0 4px 16px rgba(0,0,0,0.08);
  --shadow-lg:0 20px 60px rgba(0,0,0,0.14),0 4px 12px rgba(0,0,0,0.08);
  --sidebar-w:240px;
}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--text);font-family:'Plus Jakarta Sans',sans-serif;font-size:14px;min-height:100vh;-webkit-font-smoothing:antialiased}
.layout{display:flex;min-height:100vh}

/* ── SIDEBAR ── */
.sidebar{width:var(--sidebar-w);background:var(--sidebar-bg);display:flex;flex-direction:column;position:fixed;top:0;left:0;bottom:0;z-index:200;background-image:radial-gradient(ellipse at top left,rgba(99,102,241,0.3) 0%,transparent 60%),radial-gradient(ellipse at bottom,rgba(124,58,237,0.2) 0%,transparent 60%);transition:transform 0.28s cubic-bezier(0.4,0,0.2,1)}
.logo{padding:22px 20px 18px;border-bottom:1px solid var(--sidebar-border);display:flex;align-items:center;gap:11px}
.logo-icon{width:34px;height:34px;flex-shrink:0}
.logo-icon svg{width:100%;height:100%}
.logo-name{font-weight:800;font-size:15px;color:white;letter-spacing:0.2px;line-height:1.2}
.logo-sub{font-size:11px;color:var(--sidebar-text);margin-top:1px}
.nav{padding:16px 10px;flex:1}
.nav-label{font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:1.5px;color:rgba(255,255,255,0.3);padding:0 10px;margin-bottom:6px}
.nav-item{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:10px;cursor:pointer;font-size:13.5px;font-weight:500;color:var(--sidebar-text);transition:all 0.15s;margin-bottom:2px}
.nav-item:hover{background:rgba(255,255,255,0.08);color:white}
.nav-item.active{background:rgba(255,255,255,0.14);color:white}
.nav-icon{width:26px;height:26px;border-radius:7px;display:flex;align-items:center;justify-content:center;font-size:13px;flex-shrink:0}
.nav-item.active .nav-icon{background:rgba(255,255,255,0.18)}
.sidebar-footer{padding:12px 20px;border-top:1px solid var(--sidebar-border);display:flex;align-items:center;justify-content:space-between;gap:8px;flex-wrap:wrap}
.status-badge{display:inline-flex;align-items:center;gap:6px;font-size:11px;color:var(--sidebar-text)}
.status-dot{width:6px;height:6px;background:#34d399;border-radius:50%;box-shadow:0 0 6px #34d399;animation:pulse 2.5s ease-in-out infinite;flex-shrink:0}
@keyframes pulse{0%,100%{opacity:1;transform:scale(1)}50%{opacity:0.7;transform:scale(0.9)}}
.lang-btn{background:rgba(255,255,255,0.1);border:none;color:rgba(255,255,255,0.7);font-size:11px;font-weight:700;padding:4px 8px;border-radius:6px;cursor:pointer;font-family:'Plus Jakarta Sans',sans-serif;letter-spacing:0.5px;transition:all 0.15s}
.lang-btn:hover{background:rgba(255,255,255,0.18);color:white}
.logout-btn{background:rgba(244,63,94,0.15);border:none;color:rgba(255,180,180,0.9);font-size:11px;font-weight:600;padding:4px 9px;border-radius:6px;cursor:pointer;font-family:'Plus Jakarta Sans',sans-serif;transition:all 0.15s;display:none}
.logout-btn:hover{background:rgba(244,63,94,0.3);color:white}

/* ── MOBILE TOPBAR ── */
.mobile-topbar{display:none;position:fixed;top:0;left:0;right:0;height:52px;background:var(--sidebar-bg);z-index:150;align-items:center;padding:0 16px;gap:12px;border-bottom:1px solid var(--sidebar-border)}
.hamburger{background:none;border:none;cursor:pointer;padding:6px;border-radius:7px;display:flex;flex-direction:column;gap:4px;transition:background 0.15s}
.hamburger:hover{background:rgba(255,255,255,0.1)}
.hamburger span{display:block;width:20px;height:2px;background:white;border-radius:2px;transition:all 0.25s}
.mobile-logo-name{font-weight:800;font-size:15px;color:white;flex:1}
.sidebar-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,0.5);z-index:190;backdrop-filter:blur(2px)}
.sidebar-overlay.open{display:block}

/* ── MAIN ── */
.main{margin-left:var(--sidebar-w);flex:1;padding:32px 36px;max-width:1100px}
.page-header{display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:24px;gap:12px}
.page-title{font-size:24px;font-weight:800;color:var(--text);letter-spacing:-0.5px;line-height:1.2}
.page-sub{font-size:12.5px;color:var(--text3);margin-top:3px}
.header-actions{display:flex;gap:8px;align-items:center;flex-shrink:0}

/* ── BUTTONS ── */
.btn{display:inline-flex;align-items:center;gap:6px;padding:8px 16px;border-radius:9px;border:none;cursor:pointer;font-size:13px;font-weight:600;font-family:'Plus Jakarta Sans',sans-serif;transition:all 0.15s;white-space:nowrap;-webkit-tap-highlight-color:transparent}
.btn-primary{background:var(--accent);color:white;box-shadow:0 2px 8px rgba(79,70,229,0.3)}
.btn-primary:hover{background:#4338ca;transform:translateY(-1px);box-shadow:0 4px 14px rgba(79,70,229,0.35)}
.btn-primary:active{transform:translateY(0)}
.btn-primary:disabled{opacity:0.6;cursor:not-allowed;transform:none}
.btn-ghost{background:white;color:var(--text2);border:1.5px solid var(--border);box-shadow:var(--shadow)}
.btn-ghost:hover{background:var(--surface2);color:var(--text);border-color:var(--border2)}
.btn-danger{background:var(--red-bg);color:var(--red);border:1.5px solid #fecaca}
.btn-danger:hover{background:#fee2e2}
.btn-warning{background:#fffbeb;color:#d97706;border:1.5px solid #fde68a}
.btn-warning:hover{background:#fef3c7}
.btn-sm{padding:5px 10px;font-size:12px;border-radius:7px;gap:4px}
.btn-icon{padding:7px;border-radius:8px;width:32px;height:32px}

/* ── STATS ── */
.stats-row{display:grid;grid-template-columns:repeat(3,1fr);gap:12px;margin-bottom:22px}
.stat-card{background:white;border:1.5px solid var(--border);border-radius:14px;padding:16px 18px;position:relative;overflow:hidden;box-shadow:var(--shadow)}
.stat-card-accent{position:absolute;top:0;left:0;width:4px;height:100%}
.stat-card.total .stat-card-accent{background:linear-gradient(180deg,#4f46e5,#7c3aed)}
.stat-card.active .stat-card-accent{background:linear-gradient(180deg,#10b981,#059669)}
.stat-card.disabled .stat-card-accent{background:linear-gradient(180deg,#94a3b8,#64748b)}
.stat-label{font-size:11px;font-weight:600;color:var(--text3);text-transform:uppercase;letter-spacing:0.8px}
.stat-value{font-size:30px;font-weight:800;margin:5px 0 2px;color:var(--text);letter-spacing:-1px;line-height:1}
.stat-sub{font-size:11.5px;color:var(--text3)}

/* ── JOB LIST ── */
.section-header{display:flex;align-items:center;justify-content:space-between;margin-bottom:10px}
.section-title{font-size:14.5px;font-weight:700;color:var(--text)}
.job-list{display:flex;flex-direction:column;gap:8px}
.job-card{background:white;border:1.5px solid var(--border);border-radius:12px;padding:13px 16px;display:flex;align-items:center;gap:12px;transition:all 0.15s;box-shadow:var(--shadow);animation:slideIn 0.22s ease}
@keyframes slideIn{from{opacity:0;transform:translateY(5px)}to{opacity:1;transform:none}}
.job-card:hover{border-color:var(--border2);box-shadow:var(--shadow-md)}
.job-card.disabled{opacity:0.6}
.job-status{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.job-status.on{background:var(--green);box-shadow:0 0 6px rgba(5,150,105,0.45)}
.job-status.off{background:var(--text3)}
.job-info{flex:1;min-width:0}
.job-comment{font-size:13px;font-weight:600;color:var(--text);margin-bottom:4px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.job-command{font-family:'IBM Plex Mono',monospace;font-size:10.5px;color:#4f46e5;background:var(--accent-light);border:1px solid #c7d2fe;border-radius:4px;padding:1px 7px;display:inline-block;max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.job-schedule{font-family:'IBM Plex Mono',monospace;font-size:11px;color:#7c3aed;background:#f5f3ff;border:1px solid #ddd6fe;padding:3px 9px;border-radius:20px;white-space:nowrap;flex-shrink:0;font-weight:500}
.job-actions{display:flex;gap:5px;flex-shrink:0}

.empty-state{text-align:center;padding:56px 20px;background:white;border:1.5px dashed var(--border2);border-radius:16px}
.empty-icon{width:56px;height:56px;margin:0 auto 12px;background:var(--accent-light);border-radius:14px;display:flex;align-items:center;justify-content:center;font-size:24px}
.empty-title{font-size:14px;font-weight:600;color:var(--text2);margin-bottom:4px}
.empty-sub{font-size:12.5px;color:var(--text3)}

/* ── MODALS ── */
.modal-overlay{position:fixed;inset:0;background:rgba(15,23,42,0.55);backdrop-filter:blur(6px);z-index:1000;display:none;align-items:center;justify-content:center;padding:16px}
.modal-overlay.open{display:flex}
.modal{background:white;border:1.5px solid var(--border);border-radius:20px;width:100%;max-width:580px;max-height:94vh;overflow-y:auto;box-shadow:var(--shadow-lg);animation:modalIn 0.22s cubic-bezier(0.34,1.56,0.64,1)}
@keyframes modalIn{from{opacity:0;transform:scale(0.94) translateY(10px)}to{opacity:1;transform:none}}
.modal-sm{max-width:400px}
.modal-header{padding:18px 20px 16px;border-bottom:1.5px solid var(--border);display:flex;align-items:center;justify-content:space-between;position:sticky;top:0;background:white;border-radius:20px 20px 0 0;z-index:1}
.modal-title-wrap{display:flex;align-items:center;gap:10px}
.modal-title-icon{width:28px;height:28px;background:var(--accent-light);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:14px}
.modal-title{font-size:15px;font-weight:700;color:var(--text)}
.modal-close{background:none;border:none;cursor:pointer;color:var(--text3);font-size:16px;width:28px;height:28px;border-radius:7px;display:flex;align-items:center;justify-content:center;transition:all 0.12s}
.modal-close:hover{background:var(--surface2);color:var(--text)}
.modal-body{padding:16px 20px 18px}
.modal-footer{padding:12px 20px;border-top:1.5px solid var(--border);display:flex;gap:8px;justify-content:flex-end;align-items:center;background:var(--surface2);border-radius:0 0 20px 20px}

/* ── FORM ── */
.form-group{margin-bottom:14px}
.form-label{display:flex;align-items:center;gap:5px;font-size:11px;font-weight:700;color:var(--text2);margin-bottom:5px;text-transform:uppercase;letter-spacing:0.7px}
.form-input,.form-select,.form-textarea{width:100%;background:white;border:1.5px solid var(--border);border-radius:9px;padding:9px 12px;color:var(--text);font-size:13.5px;font-family:'Plus Jakarta Sans',sans-serif;outline:none;transition:border-color 0.15s,box-shadow 0.15s;-webkit-appearance:none}
.form-input:focus,.form-select:focus,.form-textarea:focus{border-color:var(--accent);box-shadow:0 0 0 3px rgba(79,70,229,0.1)}
.form-input::placeholder{color:var(--text3)}
.form-select{cursor:pointer;appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' fill='none'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%2394a3b8' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 12px center;padding-right:32px}
.form-textarea{min-height:110px;resize:vertical;font-family:'IBM Plex Mono',monospace;font-size:12px;line-height:1.6}
.pw-wrap{position:relative}
.pw-wrap .form-input{padding-right:42px}
.pw-toggle{position:absolute;right:10px;top:50%;transform:translateY(-50%);background:none;border:none;cursor:pointer;color:var(--text3);padding:4px;font-size:14px;transition:color 0.12s}
.pw-toggle:hover{color:var(--text2)}
.time-row{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.three-col{display:grid;grid-template-columns:1fr 1fr 1fr;gap:10px}
.choice-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:7px;margin-bottom:12px}
.choice-card{background:white;border:1.5px solid var(--border);border-radius:10px;padding:10px 8px;text-align:center;cursor:pointer;transition:all 0.15s;user-select:none;-webkit-tap-highlight-color:transparent}
.choice-card:hover{border-color:var(--border2);background:var(--surface2)}
.choice-card.selected{border-color:var(--accent);background:var(--accent-light)}
.choice-card-icon{font-size:18px;margin-bottom:3px}
.choice-card-label{font-size:11.5px;font-weight:600;color:var(--text2)}
.choice-card-sub{font-size:10px;color:var(--text3);margin-top:1px}
.choice-card.selected .choice-card-label{color:var(--accent-text)}
.choice-card.selected .choice-card-sub{color:var(--accent-text);opacity:0.7}
.cron-preview{background:var(--accent-light);border:1.5px solid #c7d2fe;border-radius:10px;padding:10px 14px;display:flex;align-items:center;gap:9px;margin-bottom:6px}
.cron-preview-label{font-size:9.5px;font-weight:700;color:var(--accent-text);text-transform:uppercase;letter-spacing:0.7px;white-space:nowrap}
.cron-preview-val{font-family:'IBM Plex Mono',monospace;font-size:13px;color:var(--accent-text);font-weight:500}
.cron-human{font-size:12px;color:var(--text2);padding:6px 11px;background:var(--surface2);border-radius:7px;border:1px solid var(--border);margin-bottom:14px}
.divider{height:1px;background:var(--border);margin:16px 0}
.section-block{margin-bottom:0}
.section-block-title{font-size:12px;font-weight:700;text-transform:uppercase;letter-spacing:0.8px;color:var(--accent-text);background:var(--accent-light);border:1px solid #c7d2fe;border-radius:8px;padding:5px 10px;margin-bottom:12px;display:inline-block}

/* ── AUTH LOGIN ── */
.login-wrap{display:flex;flex-direction:column;align-items:center;padding:8px 0 4px}
.login-logo{width:52px;height:52px;background:linear-gradient(135deg,var(--accent),var(--accent2));border-radius:14px;display:flex;align-items:center;justify-content:center;margin:0 auto 14px;box-shadow:0 4px 16px rgba(79,70,229,0.3)}
.login-logo svg{width:28px;height:28px}
.login-title{font-size:17px;font-weight:800;color:var(--text);margin-bottom:4px;text-align:center}
.login-sub{font-size:12.5px;color:var(--text3);margin-bottom:20px;text-align:center}
.login-error{font-size:12.5px;color:var(--red);background:var(--red-bg);border:1px solid #fecaca;border-radius:7px;padding:8px 12px;margin-bottom:12px;display:none}

.loading{text-align:center;padding:48px;color:var(--text3)}
.spinner{width:26px;height:26px;border:2.5px solid var(--border);border-top-color:var(--accent);border-radius:50%;animation:spin 0.7s linear infinite;margin:0 auto 12px}
@keyframes spin{to{transform:rotate(360deg)}}
.toast{position:fixed;bottom:18px;right:18px;background:white;border:1.5px solid var(--border);border-radius:12px;padding:11px 16px;font-size:13px;font-weight:500;box-shadow:var(--shadow-lg);z-index:9999;display:flex;align-items:center;gap:8px;animation:toastIn 0.25s cubic-bezier(0.34,1.56,0.64,1);max-width:300px;color:var(--text)}
@keyframes toastIn{from{opacity:0;transform:translateY(6px) scale(0.95)}to{opacity:1;transform:none}}
.toast.success{border-left:3.5px solid var(--green)}
.toast.error{border-left:3.5px solid var(--red)}

/* ── MOBILE RESPONSIVE ── */
@media(max-width:768px){
  :root{--sidebar-w:240px}
  .sidebar{transform:translateX(-100%)}
  .sidebar.mobile-open{transform:translateX(0)}
  .mobile-topbar{display:flex}
  .main{margin-left:0;padding:68px 14px 20px}
  .page-header{margin-bottom:16px}
  .page-title{font-size:20px}
  .stats-row{grid-template-columns:1fr 1fr;gap:9px;margin-bottom:16px}
  .stats-row .stat-card:last-child{grid-column:1/-1}
  .stat-value{font-size:26px}
  .job-card{flex-wrap:wrap;padding:12px 14px;gap:8px}
  .job-schedule{font-size:10.5px;padding:3px 8px}
  .job-actions{width:100%;justify-content:flex-end;gap:6px;border-top:1px solid var(--border);padding-top:8px;margin-top:2px}
  .job-actions .btn-sm{flex:1;justify-content:center}
  .modal{border-radius:16px 16px 0 0;position:fixed;bottom:0;left:0;right:0;max-width:100%;max-height:95vh;animation:sheetIn 0.28s cubic-bezier(0.4,0,0.2,1)}
  @keyframes sheetIn{from{transform:translateY(100%)}to{transform:translateY(0)}}
  .modal-overlay{align-items:flex-end;padding:0}
  .modal-sm{border-radius:16px 16px 0 0;max-width:100%}
  .modal-header{border-radius:16px 16px 0 0}
  .choice-grid{grid-template-columns:1fr 1fr}
  .time-row,.three-col{grid-template-columns:1fr 1fr}
  .header-actions .btn:not(.btn-primary) .btn-text{display:none}
  .btn-refresh-text{display:none}
}
@media(max-width:400px){
  .stats-row{grid-template-columns:1fr}
  .stats-row .stat-card:last-child{grid-column:auto}
  .job-actions .btn-sm{font-size:11px;padding:5px 6px}
}

::-webkit-scrollbar{width:5px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--border2);border-radius:3px}

/* ── LOG STYLES ── */
.log-badge{display:inline-flex;align-items:center;gap:4px;font-size:10px;font-weight:600;background:#ecfdf5;color:#059669;border:1px solid #a7f3d0;border-radius:5px;padding:2px 7px;flex-shrink:0}
.log-list{display:flex;flex-direction:column;gap:6px;max-height:320px;overflow-y:auto}
.log-item{display:flex;align-items:center;gap:10px;padding:9px 13px;background:var(--surface2);border:1.5px solid var(--border);border-radius:9px;cursor:pointer;transition:all 0.15s}
.log-item:hover{border-color:var(--accent);background:var(--accent-light)}
.log-item-time{font-size:12.5px;font-weight:600;color:var(--text);font-family:'IBM Plex Mono',monospace;flex:1}
.log-item-size{font-size:11px;color:var(--text3);flex-shrink:0}
.log-content-wrap{background:#0f172a;border-radius:10px;padding:16px;overflow:auto;max-height:440px;margin-top:2px}
.log-content{font-family:'IBM Plex Mono',monospace;font-size:12px;color:#e2e8f0;white-space:pre-wrap;word-break:break-all;line-height:1.6}
.log-empty{text-align:center;padding:40px;color:var(--text3)}
.log-empty-icon{font-size:32px;margin-bottom:10px}
.log-toolbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:10px;gap:8px;flex-wrap:wrap}
.log-filename{font-size:11.5px;color:var(--text3);font-family:'IBM Plex Mono',monospace}
.toggle-wrap{display:flex;align-items:center;gap:10px;padding:12px 14px;background:var(--surface2);border:1.5px solid var(--border);border-radius:10px;cursor:pointer;transition:all 0.15s;user-select:none}
.toggle-wrap:hover{border-color:var(--accent);background:var(--accent-light)}
.toggle-switch{width:40px;height:22px;background:#cbd5e1;border-radius:11px;position:relative;transition:background 0.2s;flex-shrink:0}
.toggle-switch.on{background:var(--accent)}
.toggle-switch::after{content:'';position:absolute;top:3px;left:3px;width:16px;height:16px;background:white;border-radius:50%;transition:transform 0.2s;box-shadow:0 1px 3px rgba(0,0,0,0.2)}
.toggle-switch.on::after{transform:translateX(18px)}
.toggle-label{font-size:13px;font-weight:500;color:var(--text2)}
.toggle-sub{font-size:11.5px;color:var(--text3);margin-top:1px}
</style>
</head>
<body>

<!-- ── LOGIN MODAL ── -->
<div class="modal-overlay" id="login-modal">
  <div class="modal modal-sm">
    <div class="modal-header" style="border-bottom:none;padding-bottom:4px">
      <div style="flex:1"></div>
    </div>
    <div class="modal-body">
      <div class="login-wrap">
        <div class="login-logo">
          <svg viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="18" cy="18" r="11" stroke="white" stroke-width="2"/>
            <line x1="18" y1="18" x2="18" y2="11" stroke="white" stroke-width="2.2" stroke-linecap="round"/>
            <line x1="18" y1="18" x2="23" y2="21" stroke="#a5b4fc" stroke-width="2.2" stroke-linecap="round"/>
            <circle cx="18" cy="18" r="2" fill="white"/>
          </svg>
        </div>
        <div class="login-title">CronPanel</div>
        <div class="login-sub" data-i18n="login_sub"></div>
        <div class="login-error" id="login-error"></div>
        <div class="form-group" style="width:100%">
          <label class="form-label" data-i18n="lbl_username"></label>
          <input type="text" class="form-input" id="login-user" data-i18n-placeholder="ph_username" autocomplete="username">
        </div>
        <div class="form-group" style="width:100%">
          <label class="form-label" data-i18n="lbl_password"></label>
          <div class="pw-wrap">
            <input type="password" class="form-input" id="login-pass" data-i18n-placeholder="ph_password" autocomplete="current-password" onkeydown="if(event.key==='Enter')doLogin()">
            <button class="pw-toggle" onclick="togglePw('login-pass',this)" type="button">👁</button>
          </div>
        </div>
        <button class="btn btn-primary" style="width:100%;justify-content:center;margin-top:4px" id="login-btn" onclick="doLogin()" data-i18n="btn_login"></button>
      </div>
    </div>
  </div>
</div>

<!-- ── LAYOUT ── -->
<div class="sidebar-overlay" id="sidebar-overlay" onclick="closeSidebar()"></div>

<div class="mobile-topbar">
  <button class="hamburger" id="hamburger" onclick="toggleSidebar()" aria-label="Menu">
    <span></span><span></span><span></span>
  </button>
  <div class="mobile-logo-name">CronPanel</div>
  <button class="lang-btn" id="lang-toggle-mobile" onclick="toggleLang()" style="font-size:11px;padding:4px 8px">EN</button>
</div>

<div class="layout">
  <nav class="sidebar" id="sidebar">
    <div class="logo">
      <div class="logo-icon">
        <svg viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect width="36" height="36" rx="9" fill="white" fill-opacity="0.15"/>
          <circle cx="18" cy="18" r="10" stroke="white" stroke-width="2"/>
          <line x1="18" y1="18" x2="18" y2="11" stroke="white" stroke-width="2.2" stroke-linecap="round"/>
          <line x1="18" y1="18" x2="23" y2="21" stroke="#a5b4fc" stroke-width="2.2" stroke-linecap="round"/>
          <circle cx="18" cy="18" r="2" fill="white"/>
          <circle cx="18" cy="8.5" r="1.2" fill="white" opacity="0.5"/>
          <circle cx="27.5" cy="18" r="1.2" fill="white" opacity="0.5"/>
          <circle cx="18" cy="27.5" r="1.2" fill="white" opacity="0.5"/>
          <circle cx="8.5" cy="18" r="1.2" fill="white" opacity="0.5"/>
        </svg>
      </div>
      <div>
        <div class="logo-name">CronPanel</div>
        <div class="logo-sub" data-i18n="sub"></div>
      </div>
    </div>
    <div class="nav">
      <div class="nav-label" data-i18n="nav_features"></div>
      <div class="nav-item active" onclick="closeSidebar()">
        <div class="nav-icon">⊞</div>
        <span data-i18n="nav_dashboard"></span>
      </div>
    </div>
    <div class="sidebar-footer">
      <div class="status-badge"><div class="status-dot"></div><span data-i18n="status_running"></span></div>
      <div style="display:flex;gap:6px;align-items:center">
        <button class="lang-btn" id="lang-toggle" onclick="toggleLang()">EN</button>
        <button class="logout-btn" id="logout-btn" onclick="doLogout()" data-i18n="btn_logout"></button>
      </div>
    </div>
  </nav>

  <main class="main" id="main-content" style="display:none">
    <div class="page-header">
      <div>
        <div class="page-title" data-i18n="dashboard_title"></div>
        <div class="page-sub" data-i18n="dashboard_sub"></div>
      </div>
      <div class="header-actions">
        <button class="btn btn-ghost btn-sm" onclick="loadJobs()">↺ <span class="btn-refresh-text" data-i18n="btn_refresh"></span></button>
        <button class="btn btn-primary btn-sm" onclick="openJobModal()">＋ <span data-i18n="btn_add"></span></button>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card total">
        <div class="stat-card-accent"></div>
        <div class="stat-label" data-i18n="stat_total_label"></div>
        <div class="stat-value" id="stat-total">—</div>
        <div class="stat-sub" data-i18n="stat_total_sub"></div>
      </div>
      <div class="stat-card active">
        <div class="stat-card-accent"></div>
        <div class="stat-label" data-i18n="stat_active_label"></div>
        <div class="stat-value" id="stat-active">—</div>
        <div class="stat-sub" data-i18n="stat_active_sub"></div>
      </div>
      <div class="stat-card disabled">
        <div class="stat-card-accent"></div>
        <div class="stat-label" data-i18n="stat_disabled_label"></div>
        <div class="stat-value" id="stat-disabled">—</div>
        <div class="stat-sub" data-i18n="stat_disabled_sub"></div>
      </div>
    </div>

    <div class="section-header">
      <div class="section-title" data-i18n="job_list_title"></div>
      <div id="last-refresh" style="font-size:11px;color:var(--text3)"></div>
    </div>
    <div id="job-list" class="job-list">
      <div class="loading"><div class="spinner"></div></div>
    </div>
  </main>
</div>

<!-- ── JOB MODAL (Single-page form) ── -->
<div class="modal-overlay" id="job-modal">
  <div class="modal">
    <div class="modal-header">
      <div class="modal-title-wrap">
        <div class="modal-title-icon" id="modal-icon">⏰</div>
        <div class="modal-title" id="modal-title"></div>
      </div>
      <button class="modal-close" onclick="closeModal()">✕</button>
    </div>
    <div class="modal-body">

      <!-- ① 任务名称 -->
      <div class="section-block">
        <div class="section-block-title" data-i18n="step_name"></div>
        <div class="form-group" style="margin-bottom:0">
          <input type="text" class="form-input" id="f-comment" data-i18n-placeholder="ph_comment" style="font-size:14px;font-weight:500;">
        </div>
      </div>

      <div class="divider"></div>

      <!-- ③ 执行时间 -->
      <div class="section-block">
        <div class="section-block-title" data-i18n="step_time"></div>
        <div class="time-row">
          <div class="form-group">
            <label class="form-label" data-i18n="lbl_hour"></label>
            <input type="number" class="form-input" id="f-hour" min="0" max="23" value="0" oninput="updatePreview()">
          </div>
          <div class="form-group">
            <label class="form-label" data-i18n="lbl_minute"></label>
            <input type="number" class="form-input" id="f-minute" min="0" max="59" value="0" oninput="updatePreview()">
          </div>
        </div>
      </div>

      <div class="divider"></div>

      <!-- ④ 执行频率 -->
      <div class="section-block">
        <div class="section-block-title" data-i18n="step_freq"></div>

        <div class="form-group">
          <label class="form-label" data-i18n="lbl_day_type"></label>
          <div class="choice-grid" id="day-type-grid">
            <div class="choice-card selected" data-val="daily" onclick="selectDayType('daily',this)">
              <div class="choice-card-icon">🔁</div>
              <div class="choice-card-label" data-i18n="opt_daily"></div>
              <div class="choice-card-sub" data-i18n="opt_daily_sub"></div>
            </div>
            <div class="choice-card" data-val="interval" onclick="selectDayType('interval',this)">
              <div class="choice-card-icon">📅</div>
              <div class="choice-card-label" data-i18n="opt_interval"></div>
              <div class="choice-card-sub" data-i18n="opt_interval_sub"></div>
            </div>
            <div class="choice-card" data-val="specific-day" onclick="selectDayType('specific-day',this)">
              <div class="choice-card-icon">📌</div>
              <div class="choice-card-label" data-i18n="opt_specific_day"></div>
              <div class="choice-card-sub" data-i18n="opt_specific_day_sub"></div>
            </div>
          </div>
        </div>
        <div id="day-interval-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_every_n_days"></label>
          <input type="number" class="form-input" id="f-days" min="1" max="365" value="2" oninput="updatePreview()">
        </div>
        <div id="day-specific-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_day_of_month"></label>
          <input type="number" class="form-input" id="f-monthday" min="1" max="31" value="1" oninput="updatePreview()">
        </div>

        <div class="form-group" style="margin-bottom:10px">
          <label class="form-label" data-i18n="lbl_month_type"></label>
          <div class="choice-grid" id="month-type-grid">
            <div class="choice-card selected" data-val="every-month" onclick="selectMonthType('every-month',this)">
              <div class="choice-card-icon">📆</div>
              <div class="choice-card-label" data-i18n="opt_every_month"></div>
              <div class="choice-card-sub" data-i18n="opt_every_month_sub"></div>
            </div>
            <div class="choice-card" data-val="month-interval" onclick="selectMonthType('month-interval',this)">
              <div class="choice-card-icon">🗓️</div>
              <div class="choice-card-label" data-i18n="opt_month_interval"></div>
              <div class="choice-card-sub" data-i18n="opt_month_interval_sub"></div>
            </div>
            <div class="choice-card" data-val="specific-month" onclick="selectMonthType('specific-month',this)">
              <div class="choice-card-icon">🎯</div>
              <div class="choice-card-label" data-i18n="opt_specific_month"></div>
              <div class="choice-card-sub" data-i18n="opt_specific_month_sub"></div>
            </div>
          </div>
        </div>
        <div id="month-interval-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_every_n_months"></label>
          <input type="number" class="form-input" id="f-month-interval" min="1" max="12" value="2" oninput="updatePreview()">
        </div>
        <div id="month-specific-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_which_month"></label>
          <select class="form-select" id="f-month" onchange="updatePreview()">
            <option value="1" data-i18n="m1"></option><option value="2" data-i18n="m2"></option>
            <option value="3" data-i18n="m3"></option><option value="4" data-i18n="m4"></option>
            <option value="5" data-i18n="m5"></option><option value="6" data-i18n="m6"></option>
            <option value="7" data-i18n="m7"></option><option value="8" data-i18n="m8"></option>
            <option value="9" data-i18n="m9"></option><option value="10" data-i18n="m10"></option>
            <option value="11" data-i18n="m11"></option><option value="12" data-i18n="m12"></option>
          </select>
        </div>

        <div class="form-group" style="margin-bottom:10px">
          <label class="form-label" data-i18n="lbl_week_type"></label>
          <div class="choice-grid" id="week-type-grid">
            <div class="choice-card selected" data-val="every-week" onclick="selectWeekType('every-week',this)">
              <div class="choice-card-icon">📅</div>
              <div class="choice-card-label" data-i18n="opt_every_week"></div>
              <div class="choice-card-sub" data-i18n="opt_every_week_sub"></div>
            </div>
            <div class="choice-card" data-val="week-interval" onclick="selectWeekType('week-interval',this)">
              <div class="choice-card-icon">🔄</div>
              <div class="choice-card-label" data-i18n="opt_week_interval"></div>
              <div class="choice-card-sub" data-i18n="opt_week_interval_sub"></div>
            </div>
            <div class="choice-card" data-val="specific-weekday" onclick="selectWeekType('specific-weekday',this)">
              <div class="choice-card-icon">📍</div>
              <div class="choice-card-label" data-i18n="opt_specific_weekday"></div>
              <div class="choice-card-sub" data-i18n="opt_specific_weekday_sub"></div>
            </div>
          </div>
        </div>
        <div id="week-interval-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_every_n_weeks"></label>
          <input type="number" class="form-input" id="f-week-interval" min="1" max="52" value="2" oninput="updatePreview()">
        </div>
        <div id="week-specific-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_which_weekday"></label>
          <select class="form-select" id="f-weekday" onchange="updatePreview()">
            <option value="0" data-i18n="wd0"></option><option value="1" data-i18n="wd1"></option>
            <option value="2" data-i18n="wd2"></option><option value="3" data-i18n="wd3"></option>
            <option value="4" data-i18n="wd4"></option><option value="5" data-i18n="wd5"></option>
            <option value="6" data-i18n="wd6"></option>
          </select>
        </div>

        <div class="form-group" style="margin-bottom:6px">
          <label style="display:flex;align-items:center;gap:8px;cursor:pointer;font-size:13px;font-weight:500;color:var(--text2)">
            <input type="checkbox" id="use-custom" onchange="toggleCustom()" style="width:14px;height:14px;accent-color:var(--accent)">
            <span data-i18n="use_custom_cron"></span>
          </label>
        </div>
        <div id="custom-cron-row" style="display:none" class="form-group">
          <label class="form-label" data-i18n="lbl_custom_cron"></label>
          <input type="text" class="form-input" id="f-custom" placeholder="* * * * *" style="font-family:'IBM Plex Mono',monospace;font-size:13px" oninput="updatePreview()">
          <div style="font-size:11px;color:var(--text3);margin-top:4px" data-i18n="custom_hint"></div>
        </div>

        <!-- Cron preview inline -->
        <div class="cron-preview">
          <div class="cron-preview-label">CRON</div>
          <div class="cron-preview-val" id="cron-preview-val">0 0 * * *</div>
        </div>
        <div class="cron-human" id="cron-human"></div>
      </div>

      <div class="divider"></div>

      <!-- ⑤ 命令 -->
      <div class="section-block">
        <div class="section-block-title" data-i18n="step_cmd"></div>
        <div class="form-group">
          <label class="form-label" data-i18n="lbl_cmd_type"></label>
          <select class="form-select" id="f-cmd-type" onchange="updateCmdType()">
            <option value="cmd" data-i18n="cmd_direct"></option>
            <option value="script-path" data-i18n="cmd_script_path"></option>
            <option value="script-content" data-i18n="cmd_script_content"></option>
          </select>
        </div>
        <div id="cmd-section-cmd">
          <div class="form-group">
            <label class="form-label" data-i18n="lbl_command"></label>
            <input type="text" class="form-input" id="f-command" style="font-family:'IBM Plex Mono',monospace;font-size:12.5px" data-i18n-placeholder="ph_command">
          </div>
        </div>
        <div id="cmd-section-script-path" style="display:none">
          <div class="form-group">
            <label class="form-label" data-i18n="lbl_script_path"></label>
            <input type="text" class="form-input" id="f-script-path" style="font-family:'IBM Plex Mono',monospace;font-size:12.5px" data-i18n-placeholder="ph_script_path">
          </div>
        </div>
        <div id="cmd-section-script-content" style="display:none">
          <div class="form-group">
            <label class="form-label" data-i18n="lbl_script_content"></label>
            <textarea class="form-textarea" id="f-script-content" data-i18n-placeholder="ph_script_content"></textarea>
            <div style="font-size:11px;color:var(--text3);margin-top:4px" data-i18n="script_note"></div>
          </div>
        </div>
      </div>

      <div class="divider"></div>

      <!-- ⑥ 日志设置 -->
      <div class="section-block">
        <div class="section-block-title" data-i18n="step_log"></div>
        <div class="toggle-wrap" onclick="toggleSaveLog()" id="log-toggle-wrap">
          <div style="flex:1">
            <div class="toggle-label" data-i18n="lbl_save_log"></div>
            <div class="toggle-sub" data-i18n="lbl_save_log_sub"></div>
          </div>
          <div class="toggle-switch" id="log-toggle-switch"></div>
        </div>
      </div>

    </div>
    <div class="modal-footer">
      <button class="btn btn-ghost btn-sm" onclick="closeModal()" data-i18n="btn_cancel"></button>
      <button class="btn btn-primary btn-sm" id="btn-submit" onclick="submitJob()" data-i18n="btn_confirm"></button>
    </div>
  </div>
</div>

<!-- ── LOG LIST MODAL ── -->
<div class="modal-overlay" id="log-list-modal">
  <div class="modal modal-sm" style="max-width:560px">
    <div class="modal-header">
      <div class="modal-title-wrap">
        <div class="modal-title-icon">📋</div>
        <div class="modal-title" id="log-list-title" data-i18n="log_list_title"></div>
      </div>
      <button class="modal-close" onclick="closeLogListModal()">✕</button>
    </div>
    <div class="modal-body">
      <div class="log-toolbar">
        <div id="log-list-sub" style="font-size:12.5px;color:var(--text3)"></div>
        <button class="btn btn-danger btn-sm" onclick="deleteAllLogs()" id="btn-delete-all-logs" data-i18n="btn_clear_logs"></button>
      </div>
      <div id="log-list-content"><div class="loading"><div class="spinner"></div></div></div>
    </div>
  </div>
</div>

<!-- ── LOG CONTENT MODAL ── -->
<div class="modal-overlay" id="log-content-modal">
  <div class="modal" style="max-width:800px">
    <div class="modal-header">
      <div class="modal-title-wrap">
        <div class="modal-title-icon">📄</div>
        <div>
          <div class="modal-title" data-i18n="log_content_title"></div>
          <div class="log-filename" id="log-content-filename"></div>
        </div>
      </div>
      <div style="display:flex;gap:8px;align-items:center">
        <button class="btn btn-danger btn-sm" onclick="deleteCurrentLog()" data-i18n="btn_delete_log"></button>
        <button class="modal-close" onclick="closeLogContentModal()">✕</button>
      </div>
    </div>
    <div class="modal-body">
      <div class="log-content-wrap">
        <pre class="log-content" id="log-content-text"></pre>
      </div>
    </div>
  </div>
</div>

<script>
// ── i18n ──────────────────────────────────────────────────
const I18N={
  zh:{
    sub:'定时任务管理器',nav_features:'功能',nav_dashboard:'仪表盘',nav_add:'添加任务',
    status_running:'服务运行中',
    dashboard_title:'仪表盘',dashboard_sub:'管理您的 Linux Crontab 定时任务',
    btn_refresh:'刷新',btn_add:'添加任务',btn_cancel:'取消',btn_confirm:'确认',btn_save:'保存',
    btn_login:'登 录',btn_logout:'退出',
    lbl_username:'用户名',lbl_password:'密码',
    ph_username:'请输入用户名',ph_password:'请输入密码',
    login_sub:'请输入凭据以继续',
    login_err_invalid:'用户名或密码错误',
    login_err_network:'网络错误，请重试',
    stat_total_label:'全部任务',stat_total_sub:'已配置的定时任务',
    stat_active_label:'运行中',stat_active_sub:'当前启用的任务',
    stat_disabled_label:'已停用',stat_disabled_sub:'已禁用的任务',
    job_list_title:'任务列表',
    refreshed_at:'更新于 ',
    empty_title:'暂无定时任务',empty_sub:'点击右上角「添加任务」开始',
    step_time:'执行时间',step_freq:'执行频率',step_cmd:'执行命令',step_name:'任务名称',
    lbl_hour:'小时 (0-23)',lbl_minute:'分钟 (0-59)',
    lbl_day_type:'每几天',lbl_month_type:'每几月',lbl_week_type:'每几周',
    opt_daily:'每天',opt_daily_sub:'每天执行',
    opt_interval:'每N天',opt_interval_sub:'指定间隔天数',
    opt_specific_day:'指定日期',opt_specific_day_sub:'每月第几天',
    opt_every_month:'每月',opt_every_month_sub:'每个月',
    opt_month_interval:'每N月',opt_month_interval_sub:'指定间隔月数',
    opt_specific_month:'指定月份',opt_specific_month_sub:'特定月份',
    opt_every_week:'每周',opt_every_week_sub:'每周执行',
    opt_week_interval:'每N周',opt_week_interval_sub:'指定间隔周数',
    opt_specific_weekday:'指定星期',opt_specific_weekday_sub:'每周特定天',
    lbl_every_n_days:'每隔几天',lbl_day_of_month:'每月第几天 (1-31)',
    lbl_every_n_months:'每隔几月',lbl_which_month:'哪个月份',
    lbl_every_n_weeks:'每隔几周',lbl_which_weekday:'星期几',
    use_custom_cron:'使用自定义 Cron 表达式（覆盖以上选项）',
    lbl_custom_cron:'Cron 表达式',custom_hint:'格式: 分 时 日 月 周，例如 0 2 * * 1',
    lbl_cmd_type:'命令类型',lbl_command:'执行命令',lbl_script_path:'脚本路径',
    lbl_script_content:'脚本内容',lbl_comment:'备注说明',optional:'(可选)',
    cmd_direct:'直接命令',cmd_script_path:'Shell 脚本路径',cmd_script_content:'编写 Shell 脚本',
    ph_command:'/usr/bin/python3 /home/user/script.py',
    ph_script_path:'/home/user/backup.sh',ph_script_content:'#!/bin/bash\necho Hello',
    ph_comment:'给这个任务起个名字，便于识别...',script_note:'脚本将保存到程序目录的 cronpanel-scripts/ 文件夹并自动执行',
    modal_add:'添加定时任务',modal_edit:'编辑定时任务',
    confirm_delete:'确认删除该任务？此操作不可撤销。',
    btn_enable:'启用',btn_disable:'停用',btn_edit:'编辑',btn_delete:'删除',
    toast_added:'任务添加成功！',toast_saved:'任务已保存',toast_deleted:'任务已删除',
    toast_err_cron:'Cron 格式错误，需要 5 个字段',toast_err_empty_cmd:'命令不能为空',
    m1:'1月',m2:'2月',m3:'3月',m4:'4月',m5:'5月',m6:'6月',
    m7:'7月',m8:'8月',m9:'9月',m10:'10月',m11:'11月',m12:'12月',
    wd0:'周日',wd1:'周一',wd2:'周二',wd3:'周三',wd4:'周四',wd5:'周五',wd6:'周六',
    human_daily:'每天 {H}:{M} 执行',human_interval:'每 {D} 天 {H}:{M} 执行',
    human_specific_day:'每月 {MD} 日 {H}:{M} 执行',human_custom:'自定义: {E}',
    step_log:'日志设置',lbl_save_log:'保存运行日志',lbl_save_log_sub:'每次执行后将输出持久化保存到文件',
    btn_view_logs:'查看日志',log_list_title:'运行日志',log_content_title:'日志内容',
    btn_clear_logs:'清空全部',btn_delete_log:'删除此条',
    log_empty_title:'暂无日志',log_empty_sub:'任务执行后日志将显示在这里',
    log_delete_confirm:'确认删除该日志文件？',log_clear_confirm:'确认清空全部日志？此操作不可撤销。',
    toast_log_deleted:'日志已删除',toast_log_cleared:'日志已清空',
  },
  en:{
    sub:'Cron Job Manager',nav_features:'Menu',nav_dashboard:'Dashboard',nav_add:'Add Job',
    status_running:'Service running',
    dashboard_title:'Dashboard',dashboard_sub:'Manage your Linux Crontab jobs',
    btn_refresh:'Refresh',btn_add:'Add Job',btn_cancel:'Cancel',btn_confirm:'Confirm',btn_save:'Save',
    btn_login:'Log In',btn_logout:'Logout',
    lbl_username:'Username',lbl_password:'Password',
    ph_username:'Enter username',ph_password:'Enter password',
    login_sub:'Enter your credentials to continue',
    login_err_invalid:'Invalid username or password',
    login_err_network:'Network error, please retry',
    stat_total_label:'Total',stat_total_sub:'Configured jobs',
    stat_active_label:'Active',stat_active_sub:'Currently enabled',
    stat_disabled_label:'Disabled',stat_disabled_sub:'Paused jobs',
    job_list_title:'Job List',
    refreshed_at:'Updated at ',
    empty_title:'No cron jobs yet',empty_sub:'Click "Add Job" to get started',
    step_time:'Execution Time',step_freq:'Frequency',step_cmd:'Command',step_name:'Task Name',
    lbl_hour:'Hour (0-23)',lbl_minute:'Minute (0-59)',
    lbl_day_type:'Day Frequency',lbl_month_type:'Month Frequency',lbl_week_type:'Week Frequency',
    opt_daily:'Every Day',opt_daily_sub:'Run daily',
    opt_interval:'Every N Days',opt_interval_sub:'Custom interval',
    opt_specific_day:'Specific Date',opt_specific_day_sub:'Day of month',
    opt_every_month:'Every Month',opt_every_month_sub:'Monthly run',
    opt_month_interval:'Every N Months',opt_month_interval_sub:'Custom interval',
    opt_specific_month:'Specific Month',opt_specific_month_sub:'Certain months',
    opt_every_week:'Every Week',opt_every_week_sub:'Weekly run',
    opt_week_interval:'Every N Weeks',opt_week_interval_sub:'Custom interval',
    opt_specific_weekday:'Specific Weekday',opt_specific_weekday_sub:'Certain days',
    lbl_every_n_days:'Every N days',lbl_day_of_month:'Day of month (1-31)',
    lbl_every_n_months:'Every N months',lbl_which_month:'Which month',
    lbl_every_n_weeks:'Every N weeks',lbl_which_weekday:'Day of week',
    use_custom_cron:'Use custom Cron expression (overrides above)',
    lbl_custom_cron:'Cron Expression',custom_hint:'Format: min hour day month weekday, e.g. 0 2 * * 1',
    lbl_cmd_type:'Command Type',lbl_command:'Command',lbl_script_path:'Script Path',
    lbl_script_content:'Script Content',lbl_comment:'Label',optional:'(optional)',
    cmd_direct:'Direct Command',cmd_script_path:'Shell Script Path',cmd_script_content:'Write Shell Script',
    ph_command:'/usr/bin/python3 /home/user/script.py',
    ph_script_path:'/home/user/backup.sh',ph_script_content:'#!/bin/bash\necho Hello',
    ph_comment:'Give this job a name...',script_note:'Script saved to cronpanel-scripts/ next to the binary',
    modal_add:'Add Cron Job',modal_edit:'Edit Cron Job',
    confirm_delete:'Delete this job? This cannot be undone.',
    btn_enable:'Enable',btn_disable:'Disable',btn_edit:'Edit',btn_delete:'Delete',
    toast_added:'Job added!',toast_saved:'Job saved',toast_deleted:'Job deleted',
    toast_err_cron:'Invalid Cron: need 5 fields',toast_err_empty_cmd:'Command cannot be empty',
    m1:'Jan',m2:'Feb',m3:'Mar',m4:'Apr',m5:'May',m6:'Jun',
    m7:'Jul',m8:'Aug',m9:'Sep',m10:'Oct',m11:'Nov',m12:'Dec',
    wd0:'Sun',wd1:'Mon',wd2:'Tue',wd3:'Wed',wd4:'Thu',wd5:'Fri',wd6:'Sat',
    human_daily:'Every day at {H}:{M}',human_interval:'Every {D} days at {H}:{M}',
    human_specific_day:'Day {MD} of each month at {H}:{M}',human_custom:'Custom: {E}',
    step_log:'Logging',lbl_save_log:'Save Run Logs',lbl_save_log_sub:'Persist stdout/stderr output to files after each run',
    btn_view_logs:'View Logs',log_list_title:'Run Logs',log_content_title:'Log Content',
    btn_clear_logs:'Clear All',btn_delete_log:'Delete',
    log_empty_title:'No logs yet',log_empty_sub:'Logs will appear here after the job runs',
    log_delete_confirm:'Delete this log file?',log_clear_confirm:'Clear all logs? This cannot be undone.',
    toast_log_deleted:'Log deleted',toast_log_cleared:'Logs cleared',
  }
};
let lang='zh';
function t(k){return(I18N[lang][k]||I18N['zh'][k]||k);}
function applyI18n(){
  document.querySelectorAll('[data-i18n]').forEach(el=>el.textContent=t(el.getAttribute('data-i18n')));
  document.querySelectorAll('[data-i18n-placeholder]').forEach(el=>el.placeholder=t(el.getAttribute('data-i18n-placeholder')));
  ['lang-toggle','lang-toggle-mobile'].forEach(id=>{const el=document.getElementById(id);if(el)el.textContent=lang==='zh'?'EN':'中文';});
  updatePreview();
}
function toggleLang(){lang=lang==='zh'?'en':'zh';applyI18n();}

// ── Session token ─────────────────────────────────────────
let sessionToken='';
function authHeaders(){return sessionToken?{'Content-Type':'application/json','Authorization':'Bearer '+sessionToken}:{'Content-Type':'application/json'};}

// ── Auth / Login ──────────────────────────────────────────
async function checkAuth(){
  try{
    const r=await fetch('/api/auth/check');const d=await r.json();
    if(d.success&&d.data){
      if(!d.data.required||d.data.loggedIn){
        if(!d.data.required) document.getElementById('logout-btn').style.display='none';
        else document.getElementById('logout-btn').style.display='';
        showApp();
      } else {
        showLoginModal();
      }
    }
  } catch(e){showLoginModal();}
}
function showLoginModal(){document.getElementById('login-modal').classList.add('open');setTimeout(()=>document.getElementById('login-user').focus(),200);}
function showApp(){
  document.getElementById('login-modal').classList.remove('open');
  document.getElementById('main-content').style.display='';
  loadJobs();
}
async function doLogin(){
  const user=document.getElementById('login-user').value.trim();
  const pass=document.getElementById('login-pass').value;
  const errEl=document.getElementById('login-error');
  const btn=document.getElementById('login-btn');
  errEl.style.display='none';
  btn.disabled=true;btn.textContent='...';
  try{
    const r=await fetch('/api/auth/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:user,password:pass})});
    const d=await r.json();
    if(d.success){
      sessionToken=d.data||'';
      document.getElementById('logout-btn').style.display='';
      showApp();
    } else {
      errEl.textContent=t('login_err_invalid');errEl.style.display='block';
    }
  } catch(e){errEl.textContent=t('login_err_network');errEl.style.display='block';}
  finally{btn.disabled=false;btn.textContent=t('btn_login');}
}
async function doLogout(){
  try{await fetch('/api/auth/logout',{method:'POST',headers:authHeaders()});}catch(e){}
  sessionToken='';
  document.getElementById('main-content').style.display='none';
  showLoginModal();
}
function togglePw(id,btn){
  const el=document.getElementById(id);
  el.type=el.type==='password'?'text':'password';
  btn.textContent=el.type==='password'?'👁':'🙈';
}

// ── Mobile sidebar ────────────────────────────────────────
function toggleSidebar(){
  const s=document.getElementById('sidebar');
  const o=document.getElementById('sidebar-overlay');
  const open=s.classList.toggle('mobile-open');
  o.classList.toggle('open',open);
}
function closeSidebar(){
  document.getElementById('sidebar').classList.remove('mobile-open');
  document.getElementById('sidebar-overlay').classList.remove('open');
}

// ── Form state ────────────────────────────────────────────
let editingId=null,dayType='daily',monthType='every-month',weekType='every-week',useCustom=false,saveLog=false;

// Detect if a command is a CronPanel-generated script path
function isManagedScript(cmd){
  // Pattern: /bin/bash <path>/cronpanel-scripts/script_<digits>.sh
  return /\/bin\/bash\s+.+\/cronpanel-scripts\/script_\d+\.sh$/.test(cmd);
}

async function openJobModal(editJob){
  editingId=editJob?editJob.id:null;
  document.getElementById('modal-title').textContent=editJob?t('modal_edit'):t('modal_add');
  document.getElementById('modal-icon').textContent=editJob?'✏️':'⏰';
  document.getElementById('btn-submit').textContent=editJob?t('btn_save'):t('btn_confirm');
  resetForm();
  if(editJob) await prefillEdit(editJob);
  document.getElementById('job-modal').classList.add('open');
}
function closeModal(){document.getElementById('job-modal').classList.remove('open');editingId=null;}
document.getElementById('job-modal').addEventListener('click',function(e){if(e.target===this)closeModal();});

async function prefillEdit(job){
  // Parse schedule
  const parts=job.schedule.split(/\s+/);
  if(parts.length===5){
    document.getElementById('f-minute').value=parts[0]==='*'?'0':parts[0];
    document.getElementById('f-hour').value=parts[1]==='*'?'0':parts[1];
    const[,, dom,mon,dow]=parts;
    if(dom.startsWith('*/')&&mon==='*'&&dow==='*'){selectDayType('interval',document.querySelector('#day-type-grid [data-val="interval"]'));document.getElementById('f-days').value=dom.replace('*/','');}
    else if(dom!=='*'&&mon==='*'&&dow==='*'){selectDayType('specific-day',document.querySelector('#day-type-grid [data-val="specific-day"]'));document.getElementById('f-monthday').value=dom;}
    if(mon.startsWith('*/')){selectMonthType('month-interval',document.querySelector('#month-type-grid [data-val="month-interval"]'));document.getElementById('f-month-interval').value=mon.replace('*/','');}
    else if(mon!=='*'){selectMonthType('specific-month',document.querySelector('#month-type-grid [data-val="specific-month"]'));document.getElementById('f-month').value=mon;}
    if(dow!=='*'){selectWeekType('specific-weekday',document.querySelector('#week-type-grid [data-val="specific-weekday"]'));document.getElementById('f-weekday').value=dow;}
  }
  document.getElementById('f-comment').value=job.comment||'';

  // Handle saveLog
  if(job.saveLog){saveLog=true;const sw=document.getElementById('log-toggle-switch');if(sw)sw.classList.add('on');}

  // Detect command type - use realCmd if available (log-wrapped jobs)
  const cmdToUse = job.realCmd || job.command;
  if(isManagedScript(cmdToUse)){
    // It's a managed script — read the file content
    document.getElementById('f-cmd-type').value='script-content';
    updateCmdType();
    try{
      const res=await fetch('/api/jobs/read-script',{method:'POST',headers:authHeaders(),body:JSON.stringify({path:cmdToUse})});
      const data=await res.json();
      if(data.success){
        document.getElementById('f-script-content').value=data.data||'';
      } else {
        document.getElementById('f-cmd-type').value='cmd';
        updateCmdType();
        document.getElementById('f-command').value=cmdToUse;
        showToast('脚本文件未找到，已显示原始命令','error');
      }
    } catch(e){
      document.getElementById('f-cmd-type').value='cmd';
      updateCmdType();
      document.getElementById('f-command').value=cmdToUse;
    }
  } else if(/^\/bin\/bash\s+/.test(cmdToUse)){
    // External script path
    document.getElementById('f-cmd-type').value='script-path';
    updateCmdType();
    document.getElementById('f-script-path').value=cmdToUse.replace(/^\/bin\/bash\s+/,'');
  } else {
    document.getElementById('f-cmd-type').value='cmd';
    updateCmdType();
    document.getElementById('f-command').value=cmdToUse;
  }
  updatePreview();
}

function selectDayType(val,el){dayType=val;document.querySelectorAll('#day-type-grid .choice-card').forEach(c=>c.classList.remove('selected'));if(el)el.classList.add('selected');document.getElementById('day-interval-row').style.display=val==='interval'?'':'none';document.getElementById('day-specific-row').style.display=val==='specific-day'?'':'none';updatePreview();}
function selectMonthType(val,el){monthType=val;document.querySelectorAll('#month-type-grid .choice-card').forEach(c=>c.classList.remove('selected'));if(el)el.classList.add('selected');document.getElementById('month-interval-row').style.display=val==='month-interval'?'':'none';document.getElementById('month-specific-row').style.display=val==='specific-month'?'':'none';updatePreview();}
function selectWeekType(val,el){weekType=val;document.querySelectorAll('#week-type-grid .choice-card').forEach(c=>c.classList.remove('selected'));if(el)el.classList.add('selected');document.getElementById('week-interval-row').style.display=val==='week-interval'?'':'none';document.getElementById('week-specific-row').style.display=val==='specific-weekday'?'':'none';updatePreview();}
function toggleCustom(){useCustom=document.getElementById('use-custom').checked;document.getElementById('custom-cron-row').style.display=useCustom?'':'none';updatePreview();}
function toggleSaveLog(){saveLog=!saveLog;const sw=document.getElementById('log-toggle-switch');if(sw)sw.classList.toggle('on',saveLog);}
function updateCmdType(){const v=document.getElementById('f-cmd-type').value;document.getElementById('cmd-section-cmd').style.display=v==='cmd'?'':'none';document.getElementById('cmd-section-script-path').style.display=v==='script-path'?'':'none';document.getElementById('cmd-section-script-content').style.display=v==='script-content'?'':'none';}

function getCronExpr(){
  if(useCustom)return(document.getElementById('f-custom').value||'').trim()||'0 0 * * *';
  const h=document.getElementById('f-hour').value||'0',m=document.getElementById('f-minute').value||'0';
  let dom='*',mon='*',dow='*';
  if(dayType==='interval'){const d=parseInt(document.getElementById('f-days').value)||2;dom='*/'+d;}
  else if(dayType==='specific-day'){dom=document.getElementById('f-monthday').value||'1';}
  if(monthType==='month-interval'){const mi=parseInt(document.getElementById('f-month-interval').value)||2;mon='*/'+mi;}
  else if(monthType==='specific-month'){mon=document.getElementById('f-month').value||'1';}
  if(weekType==='week-interval'){const wi=parseInt(document.getElementById('f-week-interval').value)||2;dow='*/'+wi;}
  else if(weekType==='specific-weekday'){dow=document.getElementById('f-weekday').value||'0';}
  return m+' '+h+' '+dom+' '+mon+' '+dow;
}
function getHumanCron(expr){
  const parts=expr.split(/\s+/);if(parts.length!==5)return expr;
  const[min,hour,dom,,]=parts;const pad=v=>String(v).padStart(2,'0');const H=pad(hour),M=pad(min);
  if(useCustom)return t('human_custom').replace('{E}',expr);
  if(dom==='*')return t('human_daily').replace('{H}',H).replace('{M}',M);
  if(dom.startsWith('*/'))return t('human_interval').replace('{D}',dom.replace('*/','')).replace('{H}',H).replace('{M}',M);
  if(dom!=='*')return t('human_specific_day').replace('{MD}',dom).replace('{H}',H).replace('{M}',M);
  return t('human_custom').replace('{E}',expr);
}
function updatePreview(){
  const expr=getCronExpr();
  const pv=document.getElementById('cron-preview-val');if(pv)pv.textContent=expr;
  const ch=document.getElementById('cron-human');if(ch)ch.textContent=getHumanCron(expr);
}

async function submitJob(){
  const ct=document.getElementById('f-cmd-type').value;
  let command='',scriptPath='',scriptContent='';
  if(ct==='cmd') command=document.getElementById('f-command').value.trim();
  else if(ct==='script-path') scriptPath=document.getElementById('f-script-path').value.trim();
  else scriptContent=document.getElementById('f-script-content').value.trim();
  // Validate
  const ok=ct==='cmd'?!!command:ct==='script-path'?!!scriptPath:!!scriptContent;
  if(!ok){showToast(t('toast_err_empty_cmd'),'error');return;}
  if(useCustom){const cron=document.getElementById('f-custom').value.trim();if(!cron||cron.split(/\s+/).length!==5){showToast(t('toast_err_cron'),'error');return;}}
  const body={mode:'custom',customCron:getCronExpr(),
    days:document.getElementById('f-days').value||'2',weekday:document.getElementById('f-weekday').value||'0',
    monthDay:document.getElementById('f-monthday').value||'1',month:document.getElementById('f-month').value||'*',
    hour:document.getElementById('f-hour').value||'0',minute:document.getElementById('f-minute').value||'0',
    comment:document.getElementById('f-comment').value.trim(),command,scriptPath,scriptContent,saveLog};
  if(editingId) body.id=editingId;
  const btn=document.getElementById('btn-submit');btn.disabled=true;
  try{
    const url=editingId?'/api/jobs/edit':'/api/jobs/add';
    const res=await fetch(url,{method:'POST',headers:authHeaders(),body:JSON.stringify(body)});
    const data=await res.json();
    if(res.status===401){showToast(t('login_err_invalid'),'error');doLogout();return;}
    if(data.success){showToast(editingId?t('toast_saved'):t('toast_added'),'success');closeModal();loadJobs();}
    else showToast(data.message||'Error','error');
  } catch(e){showToast(e.message,'error');}
  finally{btn.disabled=false;btn.textContent=editingId?t('btn_save'):t('btn_confirm');}
}

function resetForm(){
  ['f-hour','f-minute'].forEach(id=>{const el=document.getElementById(id);if(el)el.value='0';});
  ['f-days','f-month-interval','f-week-interval'].forEach(id=>{const el=document.getElementById(id);if(el)el.value='2';});
  const md=document.getElementById('f-monthday');if(md)md.value='1';
  const mo=document.getElementById('f-month');if(mo)mo.value='1';
  const wd=document.getElementById('f-weekday');if(wd)wd.value='0';
  ['f-command','f-script-path','f-script-content','f-comment','f-custom'].forEach(id=>{const el=document.getElementById(id);if(el)el.value='';});
  document.getElementById('use-custom').checked=false;document.getElementById('f-cmd-type').value='cmd';
  useCustom=false;document.getElementById('custom-cron-row').style.display='none';updateCmdType();
  saveLog=false;const sw=document.getElementById('log-toggle-switch');if(sw)sw.classList.remove('on');
  dayType='daily';monthType='every-month';weekType='every-week';
  selectDayType('daily',document.querySelector('#day-type-grid [data-val="daily"]'));
  selectMonthType('every-month',document.querySelector('#month-type-grid [data-val="every-month"]'));
  selectWeekType('every-week',document.querySelector('#week-type-grid [data-val="every-week"]'));
}

// ── Job list ──────────────────────────────────────────────
let jobs=[];
async function loadJobs(){
  try{
    const res=await fetch('/api/jobs',{headers:authHeaders()});
    if(res.status===401){doLogout();return;}
    const data=await res.json();
    if(data.success){
      jobs=data.data||[];renderJobs(jobs);updateStats(jobs);
      document.getElementById('last-refresh').textContent=t('refreshed_at')+new Date().toLocaleTimeString(lang==='zh'?'zh-CN':'en-US');
    }
  } catch(e){showToast(e.message,'error');}
}
function renderJobs(jobs){
  const el=document.getElementById('job-list');
  if(!jobs||!jobs.length){
    el.innerHTML='<div class="empty-state"><div class="empty-icon">⏰</div><div class="empty-title">'+t('empty_title')+'</div><div class="empty-sub">'+t('empty_sub')+'</div></div>';return;
  }
  el.innerHTML=jobs.map(job=>{
    const isOn=job.enabled;
    const label=job.comment||(job.command.length>48?job.command.substring(0,48)+'...':job.command);
    const displayCmd = job.saveLog&&job.realCmd ? job.realCmd : job.command;
    const id=escAttr(job.id);
    const logDirSafe=escAttr(job.logDir||'');
    const logBadge=job.saveLog?'<span class="log-badge">📋 LOG</span>':'';
    return '<div class="job-card '+(isOn?'':'disabled')+'">'+
      '<div class="job-status '+(isOn?'on':'off')+'"></div>'+
      '<div class="job-info" style="flex:1;min-width:0"><div class="job-comment" style="display:flex;align-items:center;gap:8px;flex-wrap:wrap">'+escHtml(label)+logBadge+'</div>'+
      '<span class="job-command">'+escHtml(displayCmd.length>60?displayCmd.substring(0,60)+'...':displayCmd)+'</span></div>'+
      '<span class="job-schedule">'+escHtml(job.schedule)+'</span>'+
      '<div class="job-actions">'+
        (job.saveLog?'<button class="btn btn-ghost btn-sm" onclick="openLogListModal(\''+id+'\',\''+logDirSafe+'\')">📋 '+t('btn_view_logs')+'</button>':'')+
        '<button class="btn btn-ghost btn-sm" onclick="editJob(\''+id+'\')">✏️ '+t('btn_edit')+'</button>'+
        '<button class="btn '+(isOn?'btn-warning':'btn-ghost')+' btn-sm" onclick="toggleJob(\''+id+'\')">'+
          (isOn?'⏸ '+t('btn_disable'):'▶ '+t('btn_enable'))+'</button>'+
        '<button class="btn btn-danger btn-sm" onclick="deleteJob(\''+id+'\')">✕</button>'+
      '</div></div>';
  }).join('');
}
function updateStats(jobs){
  const total=jobs.length,active=jobs.filter(j=>j.enabled).length;
  document.getElementById('stat-total').textContent=total;
  document.getElementById('stat-active').textContent=active;
  document.getElementById('stat-disabled').textContent=total-active;
}
function editJob(id){const job=jobs.find(j=>j.id===id);if(job)openJobModal(job);}
async function deleteJob(id){
  if(!confirm(t('confirm_delete')))return;
  try{
    const res=await fetch('/api/jobs/delete',{method:'POST',headers:authHeaders(),body:JSON.stringify({id})});
    if(res.status===401){doLogout();return;}
    const data=await res.json();
    if(data.success){showToast(t('toast_deleted'),'success');loadJobs();}
    else showToast(data.message,'error');
  }catch(e){showToast(e.message,'error');}
}
async function toggleJob(id){
  try{
    const res=await fetch('/api/jobs/toggle',{method:'POST',headers:authHeaders(),body:JSON.stringify({id})});
    if(res.status===401){doLogout();return;}
    const data=await res.json();
    if(data.success)loadJobs();else showToast(data.message,'error');
  }catch(e){showToast(e.message,'error');}
}

// ── Helpers ───────────────────────────────────────────────
function escHtml(s){return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');}
function escAttr(s){return String(s||'').replace(/['\"\\]/g,c=>'\\'+c);}
let toastTimer;
function showToast(msg,type){
  const ex=document.querySelector('.toast');if(ex)ex.remove();clearTimeout(toastTimer);
  const d=document.createElement('div');d.className='toast '+(type||'');
  d.innerHTML='<span>'+(type==='success'?'✓':'✕')+'</span>'+escHtml(msg);
  document.body.appendChild(d);
  toastTimer=setTimeout(()=>{d.style.opacity='0';d.style.transition='opacity 0.3s';setTimeout(()=>d.remove(),300);},3000);
}

// ── Log viewer ────────────────────────────────────────────
let currentLogDir='',currentLogFilename='';

function fmtSize(bytes){
  if(bytes<1024)return bytes+' B';
  if(bytes<1024*1024)return (bytes/1024).toFixed(1)+' KB';
  return (bytes/1024/1024).toFixed(2)+' MB';
}

async function openLogListModal(jobId,logDir){
  currentLogDir=logDir;
  document.getElementById('log-list-modal').classList.add('open');
  document.getElementById('log-list-content').innerHTML='<div class="loading"><div class="spinner"></div></div>';
  await refreshLogList();
}
function closeLogListModal(){document.getElementById('log-list-modal').classList.remove('open');currentLogDir='';}
document.getElementById('log-list-modal').addEventListener('click',function(e){if(e.target===this)closeLogListModal();});

async function refreshLogList(){
  try{
    const res=await fetch('/api/jobs/logs',{method:'POST',headers:authHeaders(),body:JSON.stringify({logDir:currentLogDir})});
    const data=await res.json();
    const logs=data.data||[];
    const sub=document.getElementById('log-list-sub');
    if(sub)sub.textContent=logs.length+' '+t('log_list_title');
    const btnDel=document.getElementById('btn-delete-all-logs');
    if(btnDel)btnDel.style.display=logs.length?'':'none';
    const cont=document.getElementById('log-list-content');
    if(!logs.length){
      cont.innerHTML='<div class="log-empty"><div class="log-empty-icon">📭</div><div style="font-weight:600;color:var(--text)">'+t('log_empty_title')+'</div><div style="font-size:12.5px;color:var(--text3);margin-top:4px">'+t('log_empty_sub')+'</div></div>';
      return;
    }
    cont.innerHTML='<div class="log-list">'+logs.map(l=>'<div class="log-item" onclick="openLogContent(\''+escAttr(l.filename)+'\')">'+
      '<span style="font-size:16px">📄</span>'+
      '<div style="flex:1;min-width:0"><div class="log-item-time">'+escHtml(l.createdAt)+'</div>'+
      '<div class="log-item-size">'+fmtSize(l.size)+'</div></div>'+
      '<button class="btn btn-danger btn-sm btn-icon" onclick="event.stopPropagation();deleteLog(\''+escAttr(l.filename)+'\')">✕</button>'+
      '</div>').join('')+'</div>';
  }catch(e){showToast(e.message,'error');}
}

async function openLogContent(filename){
  currentLogFilename=filename;
  document.getElementById('log-content-filename').textContent=filename;
  document.getElementById('log-content-text').textContent='Loading...';
  document.getElementById('log-content-modal').classList.add('open');
  try{
    const res=await fetch('/api/jobs/logs/content',{method:'POST',headers:authHeaders(),body:JSON.stringify({logDir:currentLogDir,filename})});
    const data=await res.json();
    document.getElementById('log-content-text').textContent=data.data||'(empty)';
  }catch(e){document.getElementById('log-content-text').textContent='Error: '+e.message;}
}
function closeLogContentModal(){document.getElementById('log-content-modal').classList.remove('open');currentLogFilename='';}
document.getElementById('log-content-modal').addEventListener('click',function(e){if(e.target===this)closeLogContentModal();});

async function deleteLog(filename){
  if(!confirm(t('log_delete_confirm')))return;
  try{
    await fetch('/api/jobs/logs/delete',{method:'POST',headers:authHeaders(),body:JSON.stringify({logDir:currentLogDir,filename})});
    showToast(t('toast_log_deleted'),'success');
    if(document.getElementById('log-content-modal').classList.contains('open'))closeLogContentModal();
    await refreshLogList();
  }catch(e){showToast(e.message,'error');}
}

async function deleteCurrentLog(){
  if(currentLogFilename)await deleteLog(currentLogFilename);
}

async function deleteAllLogs(){
  if(!confirm(t('log_clear_confirm')))return;
  try{
    await fetch('/api/jobs/logs/delete',{method:'POST',headers:authHeaders(),body:JSON.stringify({logDir:currentLogDir,filename:''})});
    showToast(t('toast_log_cleared'),'success');
    await refreshLogList();
  }catch(e){showToast(e.message,'error');}
}

// ── Init ──────────────────────────────────────────────────
applyI18n();
checkAuth();
setInterval(()=>{if(document.getElementById('main-content').style.display!=='none')loadJobs();},30000);
</script>
</body>
</html>`
