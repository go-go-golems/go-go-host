const express = require('express');
const ui = require('ui.dsl');
const db = require('database');

const app = express.app();

db.exec('CREATE TABLE IF NOT EXISTS visits (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL)');

app.get('/', (req, res) => {
  db.exec('INSERT INTO visits (path) VALUES (?)', req.path);
  const rows = db.query('SELECT COUNT(*) AS count FROM visits');
  return ui.page({ title: 'Hello go-go-host' },
    ui.main(
      ui.h1('Hello from go-go-host'),
      ui.p('Rendered by a hosted Goja runtime.'),
      ui.p('Visits: ' + rows[0].count)
    )
  );
});

app.get('/config-test', (req, res) => {
  try {
    db.configure('sqlite3', ':memory:');
    return res.json({ ok: false });
  } catch (err) {
    return res.json({ ok: true, message: String(err.message || err) });
  }
});
