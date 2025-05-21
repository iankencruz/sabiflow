import { Folders, House, Users, type Icon as IconType } from "@lucide/svelte";

export interface NavigationItem {
	label: string;
	href: string;
	icon?: typeof IconType, // optional icon component
	children?: NavigationItem[];
	permissions?: string[]; // optional role/permission gate
	activeMatch?: string;   // optional override for active route matching
}

export const sidebarNavigation: NavigationItem[] = [
	{
		label: 'Dashboard',
		href: '/dashboard',
		icon: House
	},
	{
		label: 'Projects',
		href: '/projects',
		icon: Folders,
		children: [
			{ label: 'All Projects', href: '/projects' },
			{ label: 'Create New', href: '/projects/new' }
		]
	},
	{
		label: 'Clients',
		href: '/clients',
		icon: Users
	},
	// {
	// 	label: 'Invoices',
	// 	href: '/invoices',
	// 	icon: 'receipt'
	// },
	// {
	// 	label: 'Settings',
	// 	href: '/settings',
	// 	icon: 'cog'
	// }
];
