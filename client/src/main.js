import './app.css'
import Login from './Login.svelte'
import Dashboard from './Dashboard.svelte'

const app = new Login({
  target: document.getElementById('app'),
})

const dashboard = new Dashboard({
  target: document.getElementById('dashboard'),
})

export default app; dashboard
