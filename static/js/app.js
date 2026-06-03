const $ = id => document.getElementById(id);
const show = id => $(id).classList.remove('hidden');
const hide = id => $(id).classList.add('hidden');

function formatCurrency(amount) {
  return 'KES ' + parseFloat(amount).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

function formatDate(dateStr) {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  return d.toLocaleDateString('en-KE', { day: '2-digit', month: 'short', year: 'numeric' });
}

function categoryBadge(cat) {
  const c = (cat || 'other').toLowerCase();
  return `<span class="badge badge-${c}">${c}</span>`;
}

function getVendorID() {
  return localStorage.getItem('vendorID') || '';
}

const loginForm = $('login-form');
if (loginForm) {
  loginForm.addEventListener('submit', async e => {
    e.preventDefault();

    const email    = $('email').value.trim();
    const password = $('password').value.trim();
    const role     = $('role').value;

    localStorage.setItem('role', role);
    localStorage.setItem('email', email);

    if (role === 'vendor') {
      localStorage.setItem('vendorID', 'vendor-001');
      window.location.href = '/vendor';
    } else {
      window.location.href = '/accountant';
    }
  });
}

const logoutBtn = $('logout-btn');
if (logoutBtn) {
  logoutBtn.addEventListener('click', () => {
    localStorage.clear();
    window.location.href = '/';
  });
}

if (window.location.pathname === '/vendor') {
  const vendorID = getVendorID();

  const nameEl = $('vendor-name');
  if (nameEl) nameEl.textContent = localStorage.getItem('email') || 'Vendor';

  loadExpenses();

  $('open-expense-form').addEventListener('click', () => show('expense-form'));
  $('cancel-expense').addEventListener('click', () => hide('expense-form'));

  $('save-expense').addEventListener('click', async () => {
    const body = {
      vendor_id:     vendorID,
      amount:        parseFloat($('exp-amount').value),
      date:          $('exp-date').value,
      category:      $('exp-category').value,
      supplier_name: $('exp-supplier').value.trim(),
      notes:         $('exp-notes').value.trim(),
    };

    if (!body.amount || !body.date) {
      alert('Please fill in at least the amount and date.');
      return;
    }

    const res = await fetch('/expenses', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify(body),
    });

    if (res.ok) {
      hide('expense-form');
      clearExpenseForm();
      loadExpenses();
    } else {
      alert('Could not save expense. Please try again.');
    }
  });

  loadInventory();

  $('open-inventory-form').addEventListener('click', () => show('inventory-form'));
  $('cancel-inventory').addEventListener('click', () => hide('inventory-form'));

  $('save-inventory').addEventListener('click', async () => {
    const body = {
      vendor_id: vendorID,
      name:      $('inv-name').value.trim(),
      quantity:  parseFloat($('inv-quantity').value),
      unit:      $('inv-unit').value.trim(),
    };

    if (!body.name || !body.quantity) {
      alert('Please fill in the item name and quantity.');
      return;
    }

    const res = await fetch('/inventory', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify(body),
    });

    if (res.ok) {
      hide('inventory-form');
      clearInventoryForm();
      loadInventory();
    } else {
      alert('Could not save item. Please try again.');
    }
  });

  loadPnL();

  $('pnl-filter-btn').addEventListener('click', () => {
    loadPnL($('pnl-from').value, $('pnl-to').value);
  });

  setupVoice();
}

if (window.location.pathname === '/accountant') {
  loadAllVendors();
  loadAllExpenses();

  $('open-vendor-form').addEventListener('click', () => show('vendor-form'));
  $('cancel-vendor').addEventListener('click', () => hide('vendor-form'));

  $('save-vendor').addEventListener('click', async () => {
    const body = {
      name:  $('v-name').value.trim(),
      email: $('v-email').value.trim(),
      role:  'vendor',
    };

    if (!body.name || !body.email) {
      alert('Please fill in name and email.');
      return;
    }

    const res = await fetch('/vendors', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify(body),
    });

    if (res.ok) {
      hide('vendor-form');
      $('v-name').value = '';
      $('v-email').value = '';
      loadAllVendors();
    } else {
      alert('Could not save vendor.');
    }
  });

  $('filter-expenses-btn').addEventListener('click', () => {
    const vendorID = $('filter-vendor').value;
    loadAllExpenses(vendorID);
  });
}

async function loadExpenses() {
  const vendorID = getVendorID();
  const res  = await fetch(`/expenses?vendorID=${vendorID}`);
  const data = await res.json();
  const tbody = $('expenses-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="6" style="text-align:center;color:var(--muted);padding:24px">No expenses yet.</td></tr>`;
    return;
  }

  tbody.innerHTML = data.map(e => `
    <tr>
      <td>${formatDate(e.date)}</td>
      <td>${e.supplier_name || '—'}</td>
      <td>${categoryBadge(e.category)}</td>
      <td><strong>${formatCurrency(e.amount)}</strong></td>
      <td>${e.notes || '—'}</td>
      <td>
        <button class="delete-btn" onclick="deleteExpense('${e.id}')">🗑</button>
      </td>
    </tr>
  `).join('');
}

async function loadInventory() {
  const vendorID = getVendorID();
  const res  = await fetch(`/inventory?vendorID=${vendorID}`);
  const data = await res.json();
  const tbody = $('inventory-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="4" style="text-align:center;color:var(--muted);padding:24px">No inventory items yet.</td></tr>`;
    return;
  }

  tbody.innerHTML = data.map(item => `
    <tr>
      <td>${item.name}</td>
      <td><strong>${item.quantity}</strong></td>
      <td>${item.unit || '—'}</td>
      <td>${formatDate(item.updated_at)}</td>
    </tr>
  `).join('');
}

async function loadPnL(from = '', to = '') {
  const vendorID = getVendorID();
  let url = `/pnl/${vendorID}`;
  if (from && to) url += `?from=${from}&to=${to}`;

  const res     = await fetch(url);
  const summary = await res.json();

  if ($('pnl-income'))   $('pnl-income').textContent   = formatCurrency(summary.total_income);
  if ($('pnl-expenses')) $('pnl-expenses').textContent = formatCurrency(summary.total_expenses);

  const profitEl = $('pnl-profit');
  if (profitEl) {
    profitEl.textContent = formatCurrency(summary.profit);
    profitEl.style.color = summary.profit >= 0 ? '#52b788' : 'var(--danger)';
  }
}

async function loadAllVendors() {
  const res  = await fetch('/vendors');
  const data = await res.json();

  const tbody = $('vendors-table-body');
  if (tbody) {
    if (!data || data.length === 0) {
      tbody.innerHTML = `<tr><td colspan="4" style="text-align:center;color:var(--muted);padding:24px">No vendors yet.</td></tr>`;
    } else {
      tbody.innerHTML = data.map(v => `
        <tr>
          <td>${v.name}</td>
          <td>${v.email}</td>
          <td>${formatDate(v.created_at)}</td>
          <td>
            <button class="btn-primary" style="padding:6px 14px;font-size:12px"
              onclick="viewVendorExpenses('${v.id}')">
              View Expenses
            </button>
          </td>
        </tr>
      `).join('');
    }
  }

  if ($('total-vendors')) $('total-vendors').textContent = data ? data.length : 0;

  const select = $('filter-vendor');
  if (select && data) {
    const options = data.map(v => `<option value="${v.id}">${v.name}</option>`).join('');
    select.innerHTML = `<option value="">All Vendors</option>` + options;
  }
}

async function loadAllExpenses(vendorID = '') {
  const url  = vendorID ? `/expenses?vendorID=${vendorID}` : '/expenses';
  const res  = await fetch(url);
  const data = await res.json();
  const tbody = $('all-expenses-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="6" style="text-align:center;color:var(--muted);padding:24px">No expenses found.</td></tr>`;
    return;
  }

  const total = data.reduce((sum, e) => sum + e.amount, 0);
  if ($('total-expenses')) $('total-expenses').textContent = formatCurrency(total);

  tbody.innerHTML = data.map(e => `
    <tr>
      <td>${formatDate(e.date)}</td>
      <td>${e.vendor_id}</td>
      <td>${e.supplier_name || '—'}</td>
      <td>${categoryBadge(e.category)}</td>
      <td><strong>${formatCurrency(e.amount)}</strong></td>
      <td>${e.notes || '—'}</td>
    </tr>
  `).join('');
}

async function deleteExpense(id) {
  if (!confirm('Delete this expense?')) return;

  const res = await fetch(`/expenses/${id}`, { method: 'DELETE' });
  if (res.ok) {
    loadExpenses();
  } else {
    alert('Could not delete expense.');
  }
}

function viewVendorExpenses(vendorID) {
  const select = $('filter-vendor');
  if (select) select.value = vendorID;
  loadAllExpenses(vendorID);
  document.getElementById('all-expenses').scrollIntoView({ behavior: 'smooth' });
}

function setupVoice() {
  const voiceBtn    = $('voice-btn');
  const voiceStatus = $('voice-status');
  if (!voiceBtn) return;

  const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
  if (!SpeechRecognition) {
    voiceBtn.textContent = '🎙 Voice not supported';
    voiceBtn.disabled = true;
    return;
  }

  const recognition = new SpeechRecognition();
  recognition.lang = 'en-KE';
  recognition.interimResults = false;

  let recording = false;

  voiceBtn.addEventListener('click', () => {
    if (!recording) {
      recognition.start();
    } else {
      recognition.stop();
    }
  });

  recognition.onstart = () => {
    recording = true;
    voiceBtn.classList.add('recording');
    voiceBtn.textContent = '🔴 Listening...';
    voiceStatus.textContent = 'Speak now — e.g. "50 shillings on transport from Mama Njeri"';
  };

  recognition.onend = () => {
    recording = false;
    voiceBtn.classList.remove('recording');
    voiceBtn.textContent = '🎙 Speak Expense';
  };

  recognition.onerror = () => {
    voiceStatus.textContent = 'Could not hear anything. Please try again.';
    voiceBtn.classList.remove('recording');
    voiceBtn.textContent = '🎙 Speak Expense';
    recording = false;
  };

  recognition.onresult = e => {
    const transcript = e.results[0][0].transcript.toLowerCase();
    voiceStatus.textContent = `Heard: "${transcript}"`;
    parseVoiceInput(transcript);
  };
}

function parseVoiceInput(text) {
  const amountMatch = text.match(/(\d+(\.\d+)?)/);
  if (amountMatch) $('exp-amount').value = amountMatch[1];

  const categories = ['food', 'transport', 'supplies', 'utilities'];
  for (const cat of categories) {
    if (text.includes(cat)) {
      $('exp-category').value = cat;
      break;
    }
  }

  const supplierMatch = text.match(/from\s+(.+)/i);
  if (supplierMatch) $('exp-supplier').value = supplierMatch[1].trim();

    if (!$('exp-date').value) {
    $('exp-date').value = new Date().toISOString().split('T')[0];
  }

   $('exp-notes').value = text;
}

function clearExpenseForm() {
  $('exp-amount').value   = '';
  $('exp-date').value     = '';
  $('exp-category').value = 'food';
  $('exp-supplier').value = '';
  $('exp-notes').value    = '';
  $('voice-status').textContent = '';
}

function clearInventoryForm() {
  $('inv-name').value     = '';
  $('inv-quantity').value = '';
  $('inv-unit').value     = '';
}



