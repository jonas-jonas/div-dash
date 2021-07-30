import ky from "ky";
import { AccountForm } from "../form/AccountForm";
import { LoginForm } from "../form/LoginForm";
import { TransactionForm } from "../form/TransactionForm";
import { Account } from "../models/account";
import { Balance } from "../models/balance";
import { SymbolChartEntry, SymbolDetails } from "../models/symbol";
import { Transaction } from "../models/transaction";
import { User } from "../models/user";

export async function getIdentity(): Promise<User> {
  const response = await ky.get("/api/auth/identity");
  return await response.json();
}

export async function postLogin(form: LoginForm): Promise<void> {
  await ky.post("/api/login", { json: form });
}

export async function getLogout(): Promise<void> {
  await ky.get("/api/auth/logout");
}

export async function getBalance(): Promise<Balance> {
  const response = await ky.get("/api/balance");
  return await response.json();
}

export async function getAccounts(): Promise<Account[]> {
  const response = await ky.get("/api/account");
  return await response.json();
}

export async function postAccount(account: AccountForm): Promise<Account> {
  const response = await ky.post("/api/account", { json: account });
  return await response.json();
}

export async function getAccount(accountId: string): Promise<Account> {
  const response = await ky.get("/api/account/" + accountId);
  return await response.json();
}

export async function getTransactions(
  accountId: string
): Promise<Transaction[]> {
  const response = await ky.get("/api/account/" + accountId + "/transaction");
  return await response.json();
}

export async function postTransaction(
  accountId: string,
  transaction: TransactionForm
): Promise<Transaction> {
  const date = new Date(transaction.date);
  const amount = parseFloat(transaction.amount);
  const price = parseFloat(transaction.price);
  const response = await ky.post("/api/account/" + accountId + "/transaction", {
    json: {
      ...transaction,
      date,
      amount,
      price,
      transactionProvider: "NONE",
    },
  });
  return await response.json();
}

export async function getSymbolDetails(
  symbolId: string
): Promise<SymbolDetails> {
  const response = await ky.get("/api/symbol/details/" + symbolId);
  return await response.json();
}

export async function getSymbolChart(
  symbolId: string
): Promise<SymbolChartEntry[]> {
  const response = await ky.get("/api/symbol/chart/" + symbolId);
  return await response.json();
}
