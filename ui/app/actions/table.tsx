import {
  Table,
  TableHead,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
  Text
} from '@tremor/react';
import Link from 'next/link';


export default function HostsTable({ hosts }: any) {
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableHeaderCell>Name</TableHeaderCell>
          <TableHeaderCell>Tags</TableHeaderCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {hosts?.map((host) => (
          <TableRow key={host.nameid}>
            <Link href={"/actions/" + host.name}>
            <TableCell>{host.name}</TableCell>
            </Link>
            <TableCell>
              <Link href={"/actions/" + host.name}>
              <Text>{host.tags}</Text>
              </Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
