import { BrowserRouter, Route, Switch } from 'react-router-dom';
import Home from "../pages/Home";
import ViewAll from "../pages/ViewAll";
import { Layout } from './Layout';

import '../styles/react-bootstrap-table2.min.css';

export function Routing() {
  return (
    <BrowserRouter>
      <Layout>
        <Switch>
          <Route exact path="/">
            <Home />
          </Route>
          <Route exact path="/all">
            <ViewAll />
          </Route>
        </Switch>
      </Layout>
    </BrowserRouter>
  )
}
