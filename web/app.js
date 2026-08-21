async function refresh() {
  const [stats, rooms, conns] = await Promise.all([
    fetch('/api/stats').then(r => r.json()),
    fetch('/api/rooms').then(r => r.json()),
    fetch('/api/conns').then(r => r.json()),
  ]);
  document.getElementById('stats').textContent = JSON.stringify(stats, null, 2);
  document.getElementById('rooms').textContent = JSON.stringify(rooms, null, 2);
  document.getElementById('conns').textContent = JSON.stringify(conns, null, 2);
}

document.getElementById('send').addEventListener('click', async () => {
  const room = document.getElementById('room').value;
  const type = document.getElementById('type').value;
  const body = document.getElementById('body').value;
  const res = await fetch('/api/broadcast', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ room, type, body }),
  });
  document.getElementById('msg').textContent = res.ok ? '已发送' : await res.text();
  refresh();
});

refresh();
setInterval(refresh, 3000);
