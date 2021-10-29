using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    /// <summary>
    /// 存储盲盒和卡牌的编号
    /// 盲盒 key: 0
    /// N卡牌 key：1
    /// E卡牌 key：2
    /// O卡牌 key：3
    /// </summary>
    public static class IndexStorage
    {
        private static readonly StorageMap IndexMap = new(Storage.CurrentContext, 0x14);

        public static BigInteger CurrentIndex(byte type)
        {
            if (type < 0 || type > 3) throw new Exception("The argument \"type\" is invalid");
            return (BigInteger)IndexMap.Get(type.ToString());
        }

        public static BigInteger NextIndex(byte type)
        {
            var value = CurrentIndex(type) + 1;
            IndexMap.Put(type.ToString(), value);
            return value;
        }

        /// <summary>
        /// 初始化NFT的编号
        /// </summary>
        public static void Initial()
        {
            IndexMap.Put("0", 0);
            IndexMap.Put("1", Neoverse.CardTypeInitialNum[0]);
            IndexMap.Put("2", Neoverse.CardTypeInitialNum[1]);
            IndexMap.Put("3", Neoverse.CardTypeInitialNum[2]);
        }
    }
}
