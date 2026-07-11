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

let currentUserID = null;
let currentUserRole = '';
let currentUserShopID = '';
let currentUserShopName = '';
let currentUserShopCode = '';

async function getSessionInfo() {
  if (currentUserID) {
    return {
      id: currentUserID,
      role: currentUserRole,
      shopID: currentUserShopID,
      shopName: currentUserShopName,
      shopCode: currentUserShopCode,
    };
  }

  const res = await fetch('/me');
  if (!res.ok) return { id: '', role: '', shopID: '', shopName: '', shopCode: '' };

  const data = await res.json();
  currentUserID = data.id || '';
  currentUserRole = data.role || '';
  currentUserShopID = data.shop_id || '';
  currentUserShopName = data.shop_name || '';
  currentUserShopCode = data.shop_code || '';
  return {
    id: currentUserID,
    role: currentUserRole,
    shopID: currentUserShopID,
    shopName: currentUserShopName,
    shopCode: currentUserShopCode,
  };
}

async function getVendorID() {
  const info = await getSessionInfo();
  return info.id;
}

const logoutBtn = $('logout-btn');
if (logoutBtn) {
  logoutBtn.addEventListener('click', () => {
    window.location.href = '/logout';
  });
}

if (window.location.pathname === '/vendor') {
  (async () => {
    const vendorID = await getVendorID();
    const sessionInfo = await getSessionInfo();

    if (sessionInfo.shopName) {
      const shopNameLabel = $('shop-name');
      if (shopNameLabel) shopNameLabel.textContent = sessionInfo.shopName;
    }

    if (sessionInfo.shopCode) {
      const shopIDValue = $('shop-id-value');
      if (shopIDValue) shopIDValue.textContent = sessionInfo.shopCode;
    }

    const copyShopIDBtn = $('copy-shop-id');
    if (copyShopIDBtn) {
      copyShopIDBtn.addEventListener('click', async () => {
        const shopID = $('shop-id-value')?.textContent || '';
        if (!shopID) return;

        try {
          if (navigator.clipboard && navigator.clipboard.writeText) {
            await navigator.clipboard.writeText(shopID);
          } else {
            const textarea = document.createElement('textarea');
            textarea.value = shopID;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);
          }

          copyShopIDBtn.textContent = 'Copied!';
          setTimeout(() => {
            copyShopIDBtn.textContent = 'Copy';
          }, 1500);
        } catch (err) {
          alert('Could not copy shop ID. Please try again.');
        }
      });
    }

    loadExpenses();
    loadInventory();
    setupVoice();

    $('open-expense-form').addEventListener('click', () => show('expense-form'));
    $('cancel-expense').addEventListener('click', () => hide('expense-form'));

    $('save-expense').addEventListener('click', async () => {
      const id = await getVendorID();
      const body = {
        vendor_id:     id,
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

    $('open-inventory-form').addEventListener('click', () => show('inventory-form'));
    $('cancel-inventory').addEventListener('click', () => hide('inventory-form'));

    $('save-inventory').addEventListener('click', async () => {
      const id = await getVendorID();
      const body = {
        vendor_id:     id,
        name:          $('inv-name').value.trim(),
        supplier_name: $('inv-supplier').value.trim(),
        status:        $('inv-status').value,
        reorder_level: parseFloat($('inv-reorder-level').value) || 0,
        expiry_date:   $('inv-expiry-date').value,
        restocked_at:  $('inv-restocked-at').value,
        quantity:      parseFloat($('inv-quantity').value),
        unit:          $('inv-unit').value.trim(),
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
  })();
}

if (window.location.pathname === '/accountant') {
  (async () => {
    const sessionInfo = await getSessionInfo();
    if (sessionInfo.shopName) {
      const shopNameLabel = $('shop-name');
      if (shopNameLabel) shopNameLabel.textContent = sessionInfo.shopName;
    }

    loadAllVendors();
    loadAllExpenses();
    loadAccountantSales();
    loadAccountantPnL();

    const openVendorForm = $('open-vendor-form');
    if (openVendorForm) {
      openVendorForm.addEventListener('click', () => show('vendor-form'));
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
    }

    $('filter-expenses-btn').addEventListener('click', () => {
      const vendorID = $('filter-vendor').value;
      loadAllExpenses(vendorID);
      loadAccountantSales(vendorID);
      loadAccountantPnL(vendorID);
    });

    const filterSalesBtn = $('filter-sales-btn');
    if (filterSalesBtn) {
      filterSalesBtn.addEventListener('click', () => {
        const vendorID = $('sales-vendor-filter')?.value || '';
        loadAccountantSales(vendorID);
        loadAccountantPnL(vendorID);
      });
    }
  })();
}

if (window.location.pathname === '/owner') {
  (async () => {
    const vendorID = await getVendorID();

    loadExpenses();
    loadInventory();
    loadSales();
    loadPnL();
    loadAllVendors();
    loadAllExpenses();
    setupVoice();

    $('open-expense-form').addEventListener('click', () => show('expense-form'));
    $('cancel-expense').addEventListener('click', () => hide('expense-form'));

    const openSaleForm = $('open-sale-form');
    if (openSaleForm) {
      openSaleForm.addEventListener('click', () => show('sale-form'));
      $('cancel-sale').addEventListener('click', () => hide('sale-form'));
      $('save-sale').addEventListener('click', async () => {
        const id = await getVendorID();
        const body = {
          vendor_id:  id,
          item_name:  $('sale-item').value.trim(),
          quantity:   parseFloat($('sale-quantity').value),
          unit_price: parseFloat($('sale-unit-price').value),
          unit_cost:  parseFloat($('sale-unit-cost').value) || 0,
          date:       $('sale-date').value,
          notes:      $('sale-notes').value.trim(),
        };

        if (!body.item_name || !body.quantity || !body.unit_price || !body.date) {
          alert('Please fill in item, quantity, price, and date.');
          return;
        }

        const res = await fetch('/sales', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body:    JSON.stringify(body),
        });

        if (res.ok) {
          hide('sale-form');
          clearSaleForm();
          loadSales();
          loadPnL($('pnl-from').value, $('pnl-to').value);
        } else {
          alert('Could not save sale. Please try again.');
        }
      });
    }

    $('save-expense').addEventListener('click', async () => {
      const id = await getVendorID();
      const body = {
        vendor_id:     id,
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

    $('open-inventory-form').addEventListener('click', () => show('inventory-form'));
    $('cancel-inventory').addEventListener('click', () => hide('inventory-form'));

    $('save-inventory').addEventListener('click', async () => {
      const id = await getVendorID();
      const body = {
        vendor_id:     id,
        name:          $('inv-name').value.trim(),
        supplier_name: $('inv-supplier').value.trim(),
        status:        $('inv-status').value,
        reorder_level: parseFloat($('inv-reorder-level').value) || 0,
        expiry_date:   $('inv-expiry-date').value,
        restocked_at:  $('inv-restocked-at').value,
        quantity:      parseFloat($('inv-quantity').value),
        unit:          $('inv-unit').value.trim(),
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

    $('pnl-filter-btn').addEventListener('click', () => {
      loadPnL($('pnl-from').value, $('pnl-to').value);
    });

    $('filter-expenses-btn').addEventListener('click', () => {
      loadAllExpenses($('filter-vendor').value);
    });
  })();
}

async function loadExpenses() {
  const vendorID = await getVendorID();
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
  const vendorID = await getVendorID();
  const res  = await fetch(`/inventory?vendorID=${vendorID}`);
  const data = await res.json();
  const tbody = $('inventory-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="8" style="text-align:center;color:var(--muted);padding:24px">No inventory items yet.</td></tr>`;
    return;
  }

  tbody.innerHTML = data.map(item => `
    <tr>
      <td>${item.name}</td>
      <td>${item.supplier_name || '—'}</td>
      <td>${item.status || '—'}</td>
      <td><strong>${item.quantity}</strong></td>
      <td>${item.reorder_level || '—'}</td>
      <td>${item.expiry_date ? formatDate(item.expiry_date) : '—'}</td>
      <td>${item.restocked_at ? formatDate(item.restocked_at) : '—'}</td>
      <td>${formatDate(item.updated_at)}</td>
    </tr>
  `).join('');
}

async function loadSales() {
  const vendorID = await getVendorID();
  const res = await fetch(`/sales?vendorID=${vendorID}`);
  const data = await res.json();
  const tbody = $('sales-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="8" style="text-align:center;color:var(--muted);padding:24px">No sales recorded yet.</td></tr>`;
    return;
  }

  tbody.innerHTML = data.map(sale => {
    const total = sale.quantity * sale.unit_price;
    const costTotal = sale.quantity * (sale.unit_cost || 0);
    const profit = total - costTotal;

    return `
      <tr>
        <td>${formatDate(sale.date)}</td>
        <td>${sale.item_name}</td>
        <td>${sale.quantity}</td>
        <td>${formatCurrency(sale.unit_price)}</td>
        <td>${formatCurrency(sale.unit_cost || 0)}</td>
        <td><strong>${formatCurrency(total)}</strong></td>
        <td>${formatCurrency(profit)}</td>
        <td><button class="delete-btn" onclick="deleteSale('${sale.id}')">🗑</button></td>
      </tr>
    `;
  }).join('');
}

async function loadPnL(from = '', to = '') {
  const vendorID = await getVendorID();
  let url = `/pnl/${vendorID}`;
  if (from && to) url += `?from=${from}&to=${to}`;

  const res     = await fetch(url);
  const summary = await res.json();

  if ($('pnl-revenue'))  $('pnl-revenue').textContent  = formatCurrency(summary.total_revenue);
  if ($('pnl-cogs'))     $('pnl-cogs').textContent     = formatCurrency(summary.total_cogs);
  if ($('pnl-expenses')) $('pnl-expenses').textContent = formatCurrency(summary.total_expenses);

  if ($('pnl-profit')) {
    $('pnl-profit').textContent = formatCurrency(summary.net_profit);
    $('pnl-profit').style.color = summary.net_profit >= 0 ? '#52b788' : 'var(--danger)';
  }
}

async function loadAccountantPnL(vendorID = '') {
  let url = '/pnl';
  if (vendorID) {
    url = `/pnl/${vendorID}`;
  }

  const res = await fetch(url);
  const summary = await res.json();

  if ($('acct-revenue'))      $('acct-revenue').textContent      = formatCurrency(summary.total_revenue);
  if ($('acct-cogs'))         $('acct-cogs').textContent         = formatCurrency(summary.total_cogs);
  if ($('acct-gross-profit')) $('acct-gross-profit').textContent = formatCurrency(summary.gross_profit);
  if ($('acct-net-profit')) {
    $('acct-net-profit').textContent = formatCurrency(summary.net_profit);
    $('acct-net-profit').style.color = summary.net_profit >= 0 ? '#52b788' : 'var(--danger)';
  }
}

async function loadAccountantSales(vendorID = '') {
  let url = '/sales';
  if (vendorID) {
    url += `?vendorID=${vendorID}`;
  }

  const res = await fetch(url);
  const data = await res.json();
  const tbody = $('sales-table-body');
  if (!tbody) return;

  if (!data || data.length === 0) {
    tbody.innerHTML = `<tr><td colspan="7" style="text-align:center;color:var(--muted);padding:24px">No sales found.</td></tr>`;
    return;
  }

  tbody.innerHTML = data.map(sale => {
    const total = sale.quantity * sale.unit_price;
    return `
      <tr>
        <td>${formatDate(sale.date)}</td>
        <td>${sale.vendor_id}</td>
        <td>${sale.item_name}</td>
        <td>${sale.quantity}</td>
        <td>${formatCurrency(sale.unit_price)}</td>
        <td>${formatCurrency(sale.unit_cost || 0)}</td>
        <td><strong>${formatCurrency(total)}</strong></td>
      </tr>
    `;
  }).join('');
}

async function deleteSale(id) {
  if (!confirm('Delete this sale?')) return;

  const res = await fetch(`/sales/${id}`, { method: 'DELETE' });
  if (res.ok) {
    loadSales();
    loadPnL($('pnl-from').value, $('pnl-to').value);
  } else {
    alert('Could not delete sale.');
  }
}

function clearSaleForm() {
  $('sale-item').value       = '';
  $('sale-quantity').value   = '';
  $('sale-unit-price').value = '';
  $('sale-unit-cost').value  = '';
  $('sale-date').value       = '';
  $('sale-notes').value      = '';
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

  const salesSelect = $('sales-vendor-filter');
  if (salesSelect && data) {
    const options = data.map(v => `<option value="${v.id}">${v.name}</option>`).join('');
    salesSelect.innerHTML = `<option value="">All Vendors</option>` + options;
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
  $('inv-name').value          = '';
  $('inv-supplier').value      = '';
  $('inv-status').value        = '';
  $('inv-reorder-level').value = '';
  $('inv-expiry-date').value   = '';
  $('inv-restocked-at').value  = '';
  $('inv-quantity').value      = '';
  $('inv-unit').value          = '';
}

window.addEventListener('pageshow', function(event) {
  if (event.persisted) {
    window.location.reload();
  }
});