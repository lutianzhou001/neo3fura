using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    public static class StateStorage
    {
        private static readonly StorageMap IndexMap = new(Storage.CurrentContext, 0x15);

        private static readonly string key = "state";

        public static void Pause() => IndexMap.Put(key, "pause");

        public static void Resume() => IndexMap.Put(key, "");

        public static string GetState() => IndexMap.Get(key) == "pause" ? "pause" : "run";

        public static bool IsPaused() => GetState() == "pause";
    }
}
