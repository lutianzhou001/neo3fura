using Neo;
using Neo.SmartContract;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Native;
using Neo.SmartContract.Framework.Services;
using System;
using System.ComponentModel;
using System.Numerics;

namespace Neoverse
{
    [DisplayName("NFT-Neoverse")]
    [SupportedStandards("NEP-11")]
    [ContractPermission("*", "onNEP11Payment")]
    [ContractTrust("0xd2a4cff31913016155e38e474a2c06d08be276cf")]
    public partial class Neoverse : Nep11Token<TokenState>
    {
        [InitialValue("NSuyiLqEfEQZsLJawCbmXW154StRUwdWoM", ContractParameterType.Hash160)]
        static readonly UInt160 Owner = default;

        private static readonly BigInteger MaxBlindBoxCount = FragmentStorage.MaxCount * 9;

        private static readonly int[] CardTypeMaxNum = { 100, 300, 3000 };

        internal static readonly int[] CardTypeInitialNum = { 0, 100, 300 };

        [Safe]
        public override string Symbol() => "N3";

        [Safe]
        public static bool IsPaused() => StateStorage.IsPaused();

        [Safe]
        public static UInt160 GetOwner()
        {
            var owner = OwnerMap.Get("owner");
            return owner != null ? (UInt160)owner : Owner;
        }

        [Safe]
        public override Map<string, object> Properties(ByteString tokenId)
        {
            StorageMap tokenMap = new(Storage.CurrentContext, Prefix_Token);
            TokenState token = (TokenState)StdLib.Deserialize(tokenMap[tokenId]);
            Map<string, object> map = new();
            map["name"] = token.Name;
            map["owner"] = token.Owner;
            map["number"] = token.Number;
            map["image"] = token.Image;
            map["video"] = token.Video;
            return map;
        }

        [Safe]
        public static TokenState GetToken(ByteString tokenId)
        {
            var TokenMap = new StorageMap(Storage.CurrentContext, Prefix_Token);
            var token = TokenMap[tokenId];
            if (token == null) throw new Exception("Neoverse::GetToken: Token \"" + tokenId + "\" does not exist.");
            return token == null ? null : (TokenState)StdLib.Deserialize(token);
        }

        /// <summary>
        /// 获得各类NFT已经发行的数量
        /// </summary>
        /// <param name="firstType">0盲盒 1碎片 2卡牌</param>
        /// <param name="secondType">盲盒：0，碎片A-I：0-8，卡牌N、E、O：0-2</param>
        /// <returns></returns>
        [Safe]
        public static BigInteger TotalMint(byte firstType, byte secondType)
        {
            if (firstType == 0)
                return IndexStorage.CurrentIndex(0);
            if (firstType == 1)
                return FragmentStorage.MaxCount - FragmentStorage.GetLastIndex(secondType); // 0~8
            if (firstType == 2)
                return IndexStorage.CurrentIndex((byte)(secondType + 1)) - CardTypeInitialNum[secondType];
            else
                throw new Exception("Neoverse::GetTokenCount: Type error.");
        }

        /// <summary>
        /// 买盲盒，2GAS一个，满20GAS，送2个
        /// </summary>
        /// <param name="from">买家地址</param>
        /// <param name="amount">购买数量</param>
        public static void OnNEP17Payment(UInt160 from, BigInteger amount, object _)
        {
            if (IsPaused()) throw new Exception("Neoverse::OnNEP17Payment: Suspension of sale.");

            amount /= 100000000;
            if (Runtime.CallingScriptHash != GAS.Hash || amount % 2 != 0)
                throw new Exception("Neoverse::OnNEP17Payment: The amount must be an integer multiple of 2GAS.");
            for (int i = 0; i < amount / 2 + amount / 20 * 2; i++)
            {
                var number = IndexStorage.NextIndex(0);
                if (number > MaxBlindBoxCount) throw new Exception("Neoverse::OnNEP17Payment: Sold out.");
                var blindBox = TokenState.CreateBlindBox(from, number);
                Mint(blindBox.Name, blindBox);
            }
        }

        /// <summary>
        /// 开盲盒
        /// </summary>
        public static bool UnBoxing(ByteString tokenId)
        {
            if (Runtime.EntryScriptHash != Runtime.CallingScriptHash) throw new Exception("Neoverse::UnBoxing: Contract calls are not allowed.");
            if (((Transaction)Runtime.ScriptContainer).Script.Length - tokenId.Length > 42) throw new Exception("Neoverse::UnBoxing: Transaction script length error.");

            return UnBoxingInternal(tokenId);
        }

        public static bool BulkUnBoxing(ByteString[] tokenList)
        {
            if (Runtime.EntryScriptHash != Runtime.CallingScriptHash) throw new Exception("Neoverse::BulkUnBoxing: Contract calls are not allowed.");

            int size = 0;
            foreach (var tokenId in tokenList)
            {
                size += tokenId.Length + 2;
                UnBoxingInternal(tokenId);
            }
            // if the size of the script without the tokenList
            // is longer than 46, it means there has some extra
            // script logic or more than two contract calls.
            // one entry call, another who knows
            if (((Transaction)Runtime.ScriptContainer).Script.Length - size > 46) throw new Exception("Neoverse::BulkUnBoxing: Transaction script length error.");

            return true;
        }

        private static bool UnBoxingInternal(ByteString tokenId)
        {
            //验证拥有者
            TokenState token = GetToken(tokenId);
            token.CheckAdmin();
            //随机数 % 9；
            //0: 碎片 A；
            //1: 碎片 B；
            //……
            //8: 碎片 I；
            //每个碎片3000个，某种碎片开没了，则给他下一个
            if (token.FirstType != 0) throw new Exception("Neoverse::UnBoxing: The token can't be unboxed");
            var mod = (byte)(Runtime.GetRandom() % 9);

            for (int i = 0; i < 9; i++)
            {
                byte fragmentIndex = ((mod + i) % 9).ToByte();

                if (FragmentStorage.GetLastIndex(fragmentIndex) > 0)
                {
                    var number = token.GetFragmentNumber(fragmentIndex);
                    var fragment = TokenState.CreateFragment(token.Owner, fragmentIndex, number);
                    if (number <= 500) CounterStorage.Increase();
                    //Burn 销毁盲盒
                    Burn(tokenId);
                    //Mint 生成碎片
                    try
                    {
                        Mint(fragment.Name, fragment);
                    }
                    catch (Exception)
                    {}
                    break;
                }
            }
            return true;
        }

        /// <summary>
        /// 碎片合成N3卡牌，一个区块即可完成
        /// </summary>
        public static bool Compound(ByteString[] tokenList)
        {
            if (Runtime.EntryScriptHash != Runtime.CallingScriptHash) throw new Exception("Neoverse::Compound: Contract calls are not allowed");

            if (tokenList.Length != 9) throw new Exception("Neoverse::Compound: The argument tokenList is invalid");
            //传入的NFT要求排好序的（减少合约执行费用），合约里对其进行验证，要求9个不同类型的碎片
            //取每个碎片的随机数 TokenState 中的 Random 字段 1~3000
            Transaction tx = (Transaction)Runtime.ScriptContainer;
            BigInteger sum = 0;
            for (int i = 0; i < 9; i++)
            {
                var fragment = GetToken(tokenList[i]);
                fragment.CheckAdmin();
                if (fragment.FirstType != 1 || fragment.SecondType != i) throw new Exception("Neoverse::Compound: Lack of fragment " + i);
                sum += fragment.Number;
            }
            CardType cardType;
            BigInteger index;
            //9个编号相加，<=4816 N卡牌，<=6411 E卡牌，其它是O卡牌
            //N卡牌限量100（1-100），E限量200（101-300），O限量2700（301-3000）
            //如果应得的卡牌发完了，则按照 N->E->O->N 的顺序发新的卡牌
            if (sum <= 4816)
            {
                (cardType, index) = GetCardType(CardType.N);
            }
            else if (sum <= 6411)
            {
                (cardType, index) = GetCardType(CardType.E);
            }
            else
            {
                (cardType, index) = GetCardType(CardType.O);
            }

            //Burn 销毁碎片
            foreach (ByteString tokenId in tokenList)
            {
                Burn(tokenId);
            }
            var card = TokenState.CreateCard(tx.Sender, (byte)cardType, index);
            //Mint 生成卡牌
            Mint(card.Name, card);
            return true;
        }

        private enum CardType
        {
            N,
            E,
            O
        }

        private static (CardType, BigInteger) GetCardType(CardType cardType)
        {
            int type = (byte)cardType;

            for (int i = 0; i < 3; i++)
            {
                var key = (type + 1).ToByte();
                BigInteger cardSurplus = IndexStorage.CurrentIndex(key);
                if (cardSurplus < CardTypeMaxNum[type])
                {
                    return ((CardType)type, IndexStorage.NextIndex(key));
                }
                else
                {
                    type = (type + 1) % 3;
                }
            }
            throw new Exception("Neoverse::GetCardType: There aren't any cards left");
        }

    }
}

