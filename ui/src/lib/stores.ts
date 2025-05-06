import type { KeycloakProfile } from "keycloak-js";
import { writable } from "svelte/store";

export const profile = writable<KeycloakProfile>()
export const token = writable<string | undefined>()
