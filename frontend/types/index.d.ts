import type { Avatar } from '#ui/types';
import * as vm from "node:vm";

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


export interface VmConfig {
  cores:number,
  memory:number,
  disk:number,
  zone:string,
}
export interface Site {
  id: number | null
  site_id: string
  name: string | null
  url: string | null
  status: string | null
  frequency: string | null
  vm_config: VmConfig
}

