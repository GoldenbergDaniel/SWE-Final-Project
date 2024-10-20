<script>
    import { Router, Link, Route } from 'svelte-routing';
    import Login from './pages/Login.svelte';
    import Dashboard from './pages/Dashboard.svelte';
    import NotFound from './pages/404.svelte';
  </script>

  <Router>
    <Route path="/" component={Login} />
    <Route path="/dashboard" component={Dashboard} />
    <Route path="*" component={NotFound} />
  </Router>