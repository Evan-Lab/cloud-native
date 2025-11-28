import { webProxyRouter } from './index.js';
import express from 'express';

const app = express();
app.use(express.json());

app.use((req, res) => {
  webProxyRouter(req, res);
});

const port = process.env.PORT || 8080;
app.listen(port, '0.0.0.0', () => {
  console.log(`Web proxy listening on port ${port}`);
});
