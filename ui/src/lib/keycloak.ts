import Keycloak from 'keycloak-js';
import { profile, token } from './stores';

export let keycloak: Keycloak

export const init = async () => {
    keycloak = new Keycloak({
        // url: 'http://localhost:8085',
        // realm: 'demo',
        // clientId: 'spa'
        url: 'https://auth.ub.unibas.ch',
        realm: 'test',
        clientId: 'graphql-demo'
    });
    try {
        await keycloak.init({
            onLoad: 'login-required',
            checkLoginIframe: false
        });
        await keycloak.loadUserProfile();
        profile.set(keycloak.profile!)
        token.set(keycloak.token)
    } catch (error) {
        console.error('Failed to initialize adapter:', error);
    }
}
