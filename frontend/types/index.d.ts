import type { Avatar } from '#ui/types';

export type UserStatus = 'subscribed' | 'unsubscribed' | 'bounced'
interface SearchParameters {
  [key: string]: any;
}

export interface User {
  name: string
  email: string
  role: 'user' | 'admin'
  avatar: Avatar
}

export interface Category {
  id: number | null
  name: string | null
  icon: string | null
  description: string | null
  parent_id: number | null
}


export interface Site {
  id: number | null
  site_id: string | null
  name: string | null
  url: string | null
  status: string | null
  no_of_crawling_per_month: number | null
  vm_config: string | null
}

