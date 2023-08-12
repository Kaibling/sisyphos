'use client';

import { Card } from '@tremor/react';
import { useState, useEffect } from 'react';
import HostsTable from './table';
import { Get } from '../lib/http';
import Link from 'next/link'

export default function HostsPage() {
  const [actions, setActions] = useState([])
  async function getActions() {
    const res = await Get("/actions");
    if (res.success){
    setActions(res.response);
  }
  }
  useEffect(() => {
    getActions();
  }, [])

  return (
    <main className="p-4 md:p-10 mx-auto max-w-7xl">
      <Link href="/actions/new">
        <button className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Create new action
        </button>
      </Link>

      <Card className="mt-6 mb-5">
       {actions.length !== 0 ? <HostsTable hosts={actions} /> : "No actions found.."}
      </Card>
    </main>
  );
}
