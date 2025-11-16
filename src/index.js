import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from "react-router-dom";

const Main = React.lazy(() => import('./Main'));
const OAuth2Handler = React.lazy(() => import('./features/schema/fields/oauth2/OAuth2Handler'));
import store from './store';
const DevToolsTheme = React.lazy(() => import('./features/theme/DevToolsTheme'));

const App = () => {
    return (
        <Provider store={store}>
            <React.Suspense fallback={<div />}> 
                <DevToolsTheme>
                    <BrowserRouter>
                        <Routes>
                            <Route exact path="/" element={<Main />} />
                            <Route path="oauth-callback" element={<OAuth2Handler />} />
                        </Routes>
                    </BrowserRouter>
                </DevToolsTheme>
            </React.Suspense>
        </Provider >
    )
}

ReactDOM.render(<App />, document.getElementById('app'));