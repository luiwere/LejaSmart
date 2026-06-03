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



