import outputJSONData from '../output/output.json'

// TODO: what if there are more islands - dynamic typing
// eslint-disable-next-line no-shadow
export enum TeamName {
    'CommonPool',
    'Team1',
    'Team2',
    'Team3',
    'Team4',
    'Team5',
    'Team6',
}

export type Transaction = {
    from: TeamName
    to: TeamName
    amount: number
}
// Custom utility type
// type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U
// export type OutputJSONType = GameStatesType & typeof outputJSONData
type ChangeKeyType<T, K extends keyof T, U> = Omit<T, K> & { [k in K]: U }

export type OutputJSONType = ChangeKeyType<
    typeof outputJSONData,
    'GameStates',
    GameStates
>

type GameStates = GameState[]
type GameState = {
    Season: number
    Turn: number
    CommonPool: number
    ClientInfos: any
    Environment: any
    DeerPopulation: any
    ForagingHistory: any
    RulesInfo: any
    IIGOHistory: IIGOHistory
    IIGORolesBudget: any
    IIGOTurnsInPower: any
    IIGOTaxAmount: any
    IIGOAllocationMap: any
    IIGOSanctionMap: any
    IIGOSanctionCache: any
    IIGOHistoryCache: any
    IIGORoleMonitoringCache: any
    IITOTransactions: IITOTransactions
    SpeakerID: string
    JudgeID: string
    PresidentID: string
}

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: string]: Accountability[] | undefined
}

export type Accountability = {
    ClientID: string
    Pairs: VariableValuePair[]
}

export type VariableValuePair = {
    VariableName: string
    Values: number[]
}

// IITO Types
export type IITOTransactions = {
    [key: number]: GiftResponseDict
}

export type GiftResponseDict = {
    [team: number]: GiftResponse
}

export type GiftResponse = {
    AcceptedAmount: number
    Reason: number
}
