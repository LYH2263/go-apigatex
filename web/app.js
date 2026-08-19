async function refresh() {
  const routes = await (await fetch('/api/routes')).json();
  const el = document.getElementById('routes');
  el.innerHTML = routes.map(r =>
    `<div class="route"><span>${r.ID} ${r.Path} → ${r.Upstream}</span>` +
    `<button data-id="${r.ID}" class="del">删除</button></div>`
  ).join('') || '<p>暂无路由</p>';
  el.querySelectorAll('.del').forEach(btn => btn.onclick = async () => {
    await fetch('/api/routes/' + btn.dataset.id, { method: 'DELETE' });
    refresh();
  });
  document.getElementById('stats').textContent = JSON.stringify(await (await fetch('/api/stats')).json(), null, 2);
}
document.getElementById('add-form').onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  await fetch('/api/routes', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      ID: fd.get('id'), Path: fd.get('path'), Kind: 'prefix',
      Upstream: fd.get('upstream'), Auth: !!fd.get('auth'), Methods: ['GET','POST','PUT','DELETE']
    })
  });
  e.target.reset();
  refresh();
};
document.getElementById('try-form').onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  const res = await fetch('/api/try', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ method: fd.get('method'), path: fd.get('path'), body: fd.get('body') || '' })
  });
  document.getElementById('try-out').textContent = await res.text();
};
refresh();
