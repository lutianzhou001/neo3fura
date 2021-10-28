using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Native;
using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    public static class FragmentStorage
    {
        private static readonly StorageMap LastIndexMap = new(Storage.CurrentContext, 0x12);
        private static readonly StorageMap RandomMap = new(Storage.CurrentContext, 0x13);

        public const int MaxCount = 3000;

        public static void Initial()
        {
            for (int i = 0; i <= 8; i++)
            {
                LastIndexMap.Put(i.ToString(), MaxCount);
            }
        }

        public static BigInteger GetLastIndex(byte type)
        {
            if (type < 0 || type > 8) throw new Exception("Neoverse::LastIndexStorage: Unkown type");
            return (BigInteger)LastIndexMap.Get(type.ToString());
        }

        private static void ReduceLastIndex(byte type)
        {
            BigInteger oldIndexNumber = GetLastIndex(type);
            if (oldIndexNumber <= 0) throw new Exception("Neoverse::LastIndexStorage: Sold out");
            LastIndexMap.Put(type.ToString(), oldIndexNumber - 1);
        }

        public static BigInteger GetFragmentNumber(this TokenState token, byte type)
        {
            if (type < 0 || type > 8) throw new Exception("Neoverse::GetFragmentNumber: Unknown Type");

            BigInteger lastIndex = GetLastIndex(type);
            var random = Runtime.GetRandom() % lastIndex + 1;
            var randomKey = GetKey(type, random);
            var randomValue = RandomMap.Get(randomKey);

            if (random <= 500 && randomValue == null)
            {
                var count = CounterStorage.Get();
                if (count >= 100)
                {
                    var m = count / 100 * 10;
                    m = lastIndex - random < m ? lastIndex - random : m;
                    var offset = random % m;
                    random += offset;
                    randomKey = GetKey(type, random);
                    randomValue = RandomMap.Get(randomKey);
                }
            }

            BigInteger result = (randomValue == null) ? random : (BigInteger)randomValue;

            var lastIndexKey = GetKey(type, lastIndex);
            var lastValue = RandomMap.Get(lastIndexKey);
            RandomMap.Put(randomKey, lastValue == null ? GetLastIndex(type) : (BigInteger)lastValue);

            ReduceLastIndex(type);

            return result;
        }

        private static byte[] GetKey(byte type, BigInteger random) => new byte[] { type }.Concat(random.ToByteArray());
    }
}
