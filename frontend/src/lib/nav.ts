import { House, Volleyball, Telescope, Network, Trophy } from '@lucide/svelte';
import type { Component } from 'svelte';

export interface NavItem {
	href: string;
	label: string;
	icon: Component;
}

export const navItems: NavItem[] = [
	{ href: '/', label: 'Inicio', icon: House },
	{ href: '/tips', label: 'Tips', icon: Volleyball },
	{ href: '/forecast', label: 'Pronóstico', icon: Telescope },
	{ href: '/tournament', label: 'Cuadro', icon: Network },
	{ href: '/leagues', label: 'Ligas', icon: Trophy }
];

export function isActive(href: string, path: string): boolean {
	return href === '/' ? path === '/' : path.startsWith(href);
}
