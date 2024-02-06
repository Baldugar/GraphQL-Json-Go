export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Map: { input: { [key: string]: any }; output: { [key: string]: any }; }
};

export type Model = {
  __typename?: 'Model';
  fields: Array<ModelField>;
  name: Scalars['String']['output'];
};

export type ModelField = {
  __typename?: 'ModelField';
  isNullable: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  subFields?: Maybe<Array<ModelField>>;
  type: ModelFieldEnum;
};

export enum ModelFieldEnum {
  ARRAY = 'ARRAY',
  BOOLEAN = 'BOOLEAN',
  FLOAT = 'FLOAT',
  INT = 'INT',
  OBJECT = 'OBJECT',
  STRING = 'STRING'
}

export type ModelFieldInput = {
  isNullable: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  subFields?: InputMaybe<Array<ModelFieldInput>>;
  type: ModelFieldEnum;
};

export type Mutation = {
  __typename?: 'Mutation';
  createModel: Model;
};


export type MutationcreateModelArgs = {
  fields: Array<ModelFieldInput>;
  name: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  getModels: Array<Model>;
  sendInformaton: Scalars['Boolean']['output'];
};


export type QuerysendInformatonArgs = {
  info: Scalars['Map']['input'];
  modelName: Scalars['String']['input'];
};
