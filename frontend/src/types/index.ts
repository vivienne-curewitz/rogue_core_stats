/**
 * Branded type for UUIDs to ensure type safety.
 */
export type UUID = string & { readonly __brand: unique symbol };

/**
 * Helper to cast or parse a string to a UUID.
 */
export function asUUID(value: string): UUID {
  return value as UUID;
}

/**
 * Helper to check if a string is a valid UUID.
 */
export function isUUID(value: string): value is UUID {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(value);
}

/**
 * Branded type for Integers.
 */
export type Int = number & { readonly __brand: unique symbol };

/**
 * Helper to convert a number or string to a safe Int.
 */
export function asInt(value: number | string): Int {
  const num = typeof value === 'string' ? parseInt(value, 10) : Math.floor(value);
  if (isNaN(num)) {
    return 0 as Int;
  }
  return num as Int;
}

/**
 * Converts a base64 string to a Uint8Array.
 */
export function base64ToBytes(base64: string): Uint8Array {
  const binString = atob(base64);
  return Uint8Array.from(binString, (m) => m.codePointAt(0)!);
}

/**
 * Converts a Uint8Array to a base64 string.
 */
export function bytesToBase64(bytes: Uint8Array): string {
  const binString = Array.from(bytes, (x) => String.fromCodePoint(x)).join("");
  return btoa(binString);
}

export interface Upgrade {
  RunId: string;
  PlayerId: string;
  UpgradeId: string;
  Quantity: Int;
  Reference: string; // reference UUID string or link
}

export interface Item {
  RunId: string;
  PlayerId: string;
  ItemId: string;
  Reference: string; // reference UUID string
}

export interface RunStatus {
  RunId: string;
  Status: boolean;
}

export interface RunOverview {
  RunId: string;
  PlayerId: string;
  Status: boolean;
  BossId: string;
  Depth: Int;
  CharacterId: string;
  PlayerDamage: number;
  OverkillDamage: number;
  PlayerKills: Int;
  PlayerDeaths: Int;
  TotalStages: Int;
  CompletedStages: Int;
  Runtime: Int;
  PlayerRank: Int;
  CharacterRank: Int;
  CharacterStars: Int;
  MineralsMined: number;
  MaxArmor: number;
  MaxHealth: number;
  HealthRestored: number;
  Timestamp: Int;
}

export enum UploadStatus {
  Pending = 0,
  InProgress = 1,
  Completed = 2,
  Failed = 3,
}

export const UploadStatusStrings = {
  [UploadStatus.Pending]: 'Pending',
  [UploadStatus.InProgress]: 'In Progress',
  [UploadStatus.Completed]: 'Completed',
  [UploadStatus.Failed]: 'Failed',
} as const;

export type UploadStatusString = typeof UploadStatusStrings[UploadStatus] | 'Unknown';

export function parseUploadStatus(val: number | string): UploadStatus {
  if (typeof val === 'number') {
    return val in UploadStatus ? (val as UploadStatus) : UploadStatus.Pending;
  }
  switch (val.trim()) {
    case 'Pending':
      return UploadStatus.Pending;
    case 'In Progress':
      return UploadStatus.InProgress;
    case 'Completed':
      return UploadStatus.Completed;
    case 'Failed':
      return UploadStatus.Failed;
    default: {
      const parsed = parseInt(val, 10);
      if (!isNaN(parsed) && parsed in UploadStatus) {
        return parsed as UploadStatus;
      }
      return UploadStatus.Pending;
    }
  }
}

export interface SaveDataTask {
  Data: Uint8Array | string; // Uint8Array or base64 encoded string
  ID: UUID;
}

/**
 * Converts a raw object to a typed Upgrade, ensuring ints and uuids are parsed.
 */
export function parseUpgrade(raw: Record<string, unknown> | null | undefined): Upgrade {
  const r = raw ?? {};
  return {
    RunId: String(r.RunId ?? ''),
    PlayerId: String(r.PlayerId ?? ''),
    UpgradeId: String(r.UpgradeId ?? ''),
    Quantity: asInt((r.Quantity ?? 0) as number | string),
    Reference: String(r.Reference ?? ''),
  };
}

/**
 * Converts a raw object to a typed Item, ensuring uuids are parsed.
 */
export function parseItem(raw: Record<string, unknown> | null | undefined): Item {
  const r = raw ?? {};
  return {
    RunId: String(r.RunId ?? ''),
    PlayerId: String(r.PlayerId ?? ''),
    ItemId: String(r.ItemId ?? ''),
    Reference: String(r.Reference ?? ''),
  };
}

/**
 * Converts a raw object to a typed RunStatus.
 */
export function parseRunStatus(raw: Record<string, unknown> | null | undefined): RunStatus {
  const r = raw ?? {};
  return {
    RunId: String(r.RunId ?? ''),
    Status: Boolean(r.Status),
  };
}

/**
 * Converts a raw object to a typed RunOverview, ensuring ints are parsed.
 */
export function parseRunOverview(raw: Record<string, unknown> | null | undefined): RunOverview {
  const r = raw ?? {};
  return {
    RunId: String(r.RunId ?? ''),
    PlayerId: String(r.PlayerId ?? ''),
    Status: Boolean(r.Status),
    BossId: String(r.BossId ?? ''),
    Depth: asInt((r.Depth ?? 0) as number | string),
    CharacterId: String(r.CharacterId ?? ''),
    PlayerDamage: Number(r.PlayerDamage ?? 0),
    OverkillDamage: Number(r.OverkillDamage ?? 0),
    PlayerKills: asInt((r.PlayerKills ?? 0) as number | string),
    PlayerDeaths: asInt((r.PlayerDeaths ?? 0) as number | string),
    TotalStages: asInt((r.TotalStages ?? 0) as number | string),
    CompletedStages: asInt((r.CompletedStages ?? 0) as number | string),
    Runtime: asInt((r.Runtime ?? 0) as number | string),
    PlayerRank: asInt((r.PlayerRank ?? 0) as number | string),
    CharacterRank: asInt((r.CharacterRank ?? 0) as number | string),
    CharacterStars: asInt((r.CharacterStars ?? 0) as number | string),
    MineralsMined: Number(r.MineralsMined ?? 0),
    MaxArmor: Number(r.MaxArmor ?? 0),
    MaxHealth: Number(r.MaxHealth ?? 0),
    HealthRestored: Number(r.HealthRestored ?? 0),
    Timestamp: asInt((r.Timestamp ?? 0) as number | string),
  };
}

/**
 * Converts a raw object to a typed SaveDataTask, ensuring UUID is parsed.
 */
export function parseSaveDataTask(raw: Record<string, unknown> | null | undefined): SaveDataTask {
  const r = raw ?? {};
  let data: Uint8Array | string = (r.Data ?? '') as Uint8Array | string;
  if (typeof data === 'string' && data) {
    try {
      data = base64ToBytes(data);
    } catch {
      // keep as string if not valid base64
    }
  }
  return {
    Data: data,
    ID: asUUID(String(r.ID ?? '')),
  };
}

