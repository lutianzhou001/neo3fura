using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    /// <summary>
    /// 碎片编号小于等于500的数量，用来计算偏移量，以保证开出小碎片越来越难
    /// </summary>
    public static class CounterStorage
    {
        private static readonly StorageMap CounterMap = new(Storage.CurrentContext, 0x11);

        private static readonly string key = "counter";

        public static void Initial()
        {
            CounterMap.Put(key, 0);
        }
        public static BigInteger Get() => (BigInteger)CounterMap.Get(key);

        public static void Increase() => CounterMap.Put(key, Get() + 1);
    }
}
